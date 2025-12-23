package rabbitmq

import (
	"context"
	"time"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog/log"

	goqueue "github.com/turtlepavlo/async-queue-service"
	"github.com/turtlepavlo/async-queue-service/errors"
	headerKey "github.com/turtlepavlo/async-queue-service/headers/key"
	headerVal "github.com/turtlepavlo/async-queue-service/headers/value"
	"github.com/turtlepavlo/async-queue-service/interfaces"
	"github.com/turtlepavlo/async-queue-service/internal/publisher"
	"github.com/turtlepavlo/async-queue-service/middleware"
	publisherOpts "github.com/turtlepavlo/async-queue-service/options/publisher"
)

const (
	DefaultChannelPoolSize = 5
)

type rabbitMQ struct {
	channelPool *ChannelPool
	option      *publisherOpts.PublisherOption
}

func NewPublisher(
	opts ...publisherOpts.PublisherOptionFunc,
) publisher.Publisher {
	opt := publisherOpts.DefaultPublisherOption()
	for _, o := range opts {
		o(opt)
	}

	conn := opt.RabbitMQPublisherConfig.Conn
	channelPool := NewChannelPool(conn, opt.RabbitMQPublisherConfig.PublisherChannelPoolSize)
	ch, err := channelPool.Get()
	if err != nil {
		log.Fatal().Err(err).Msg("error getting channel from pool")
	}
	defer channelPool.Return(ch)

	return &rabbitMQ{
		channelPool: channelPool,
		option:      opt,
	}
}

func (r *rabbitMQ) Publish(ctx context.Context, m interfaces.Message) (err error) {
	if m.ContentType == "" {
		m.ContentType = publisherOpts.DefaultContentType
	}
	publishFunc := middleware.ApplyPublisherMiddleware(
		r.buildPublisher(),
		r.option.Middlewares...,
	)
	return publishFunc(ctx, m)
}

func (r *rabbitMQ) buildPublisher() interfaces.PublisherFunc {
	return func(ctx context.Context, m interfaces.Message) (err error) {
		id := m.ID
		if id == "" {
			id = uuid.New().String()
		}

		timestamp := m.Timestamp
		if timestamp.IsZero() {
			timestamp = time.Now()
		}

		defaultHeaders := map[string]any{
			headerKey.AppID:              r.option.PublisherID,
			headerKey.MessageID:          id,
			headerKey.PublishedTimestamp: timestamp.Format(time.RFC3339),
			headerKey.RetryCount:         0,
			headerKey.ContentType:        string(m.ContentType),
			headerKey.QueueServiceAgent:  string(headerVal.RabbitMQ),
			headerKey.SchemaVer:          headerVal.GoquMessageSchemaVersionV1,
		}

		headers := amqp.Table{}
		for key, value := range defaultHeaders {
			headers[key] = value
		}
		for key, value := range m.Headers {
			headers[key] = value
		}

		m.Headers = headers
		m.ServiceAgent = headerVal.RabbitMQ
		m.Timestamp = timestamp
		m.ID = id
		encoder, ok := goqueue.GetGoQueueEncoding(m.ContentType)
		if !ok {
			return errors.ErrEncodingFormatNotSupported
		}

		data, err := encoder.Encode(ctx, m)
		if err != nil {
			return err
		}

		ch, err := r.channelPool.Get()
		if err != nil {
			return err
		}
		defer r.channelPool.Return(ch)
		return ch.PublishWithContext(
			ctx,
			m.Topic,  // exchange
			m.Action, // routing-key
			false,    // mandatory
			false,    // immediate
			amqp.Publishing{
				Headers:     headers,
				ContentType: string(m.ContentType),
				Body:        data,
				Timestamp:   timestamp,
				AppId:       r.option.PublisherID,
			},
		)
	}
}

func (r *rabbitMQ) Close(_ context.Context) (err error) {
	return r.channelPool.Close()
}
