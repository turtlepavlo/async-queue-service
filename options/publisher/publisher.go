package publisher

import (
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"

	headerVal "github.com/turtlepavlo/async-queue-service/headers/value"
	"github.com/turtlepavlo/async-queue-service/interfaces"
	"github.com/turtlepavlo/async-queue-service/options"
)

const DefaultContentType = headerVal.ContentTypeJSON

type PublisherOption struct {
	PublisherID             string
	Middlewares             []interfaces.PublisherMiddlewareFunc
	RabbitMQPublisherConfig *RabbitMQPublisherConfig
}

type PublisherOptionFunc func(opt *PublisherOption)

func WithPublisherID(id string) PublisherOptionFunc {
	return func(opt *PublisherOption) { opt.PublisherID = id }
}

func WithMiddlewares(middlewares ...interfaces.PublisherMiddlewareFunc) PublisherOptionFunc {
	return func(opt *PublisherOption) { opt.Middlewares = middlewares }
}

var DefaultPublisherOption = func() *PublisherOption {
	return &PublisherOption{
		Middlewares: []interfaces.PublisherMiddlewareFunc{},
		PublisherID: uuid.New().String(),
	}
}

type RabbitMQPublisherConfig struct {
	PublisherChannelPoolSize int
	Conn                     *amqp.Connection
}

func WithRabbitMQPublisherConfig(rabbitMQOption *RabbitMQPublisherConfig) PublisherOptionFunc {
	return func(opt *PublisherOption) { opt.RabbitMQPublisherConfig = rabbitMQOption }
}

const (
	PublisherPlatformRabbitMQ     = options.PlatformRabbitMQ
	PublisherPlatformGooglePubSub = options.PlatformGooglePubSub
	PublisherPlatformSNS          = options.PlatformSNS
)
