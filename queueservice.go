package goqueue

import (
	"context"
	"errors"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/turtlepavlo/async-queue-service/interfaces"
	"github.com/turtlepavlo/async-queue-service/internal/consumer"
	"github.com/turtlepavlo/async-queue-service/internal/publisher"
	"github.com/turtlepavlo/async-queue-service/options"
)

// QueueService is the top-level entry point for consuming and publishing messages.
type QueueService struct {
	consumer         consumer.Consumer
	handler          interfaces.InboundMessageHandler
	publisher        publisher.Publisher
	NumberOfConsumer int
}

func NewQueueService(opts ...options.GoQueueOptionFunc) *QueueService {
	opt := options.DefaultGoQueueOption()
	for _, o := range opts {
		o(opt)
	}
	return &QueueService{
		consumer:         opt.Consumer,
		handler:          opt.MessageHandler,
		publisher:        opt.Publisher,
		NumberOfConsumer: opt.NumberOfConsumer,
	}
}

// Start spawns NumberOfConsumer goroutines, each running the consumer loop.
// Blocks until the context is cancelled or all consumers exit.
func (qs *QueueService) Start(ctx context.Context) (err error) {
	if qs.consumer == nil {
		return errors.New("consumer is not defined")
	}
	if qs.handler == nil {
		return errors.New("handler is not defined")
	}

	g, ctx := errgroup.WithContext(ctx)
	for i := range qs.NumberOfConsumer {
		meta := map[string]any{
			"consumer_id":  i,
			"started_time": time.Now(),
		}
		g.Go(func() error {
			return qs.consumer.Consume(ctx, qs.handler, meta)
		})
	}

	return g.Wait()
}

func (qs *QueueService) Stop(ctx context.Context) error {
	if qs.consumer != nil {
		if err := qs.consumer.Stop(ctx); err != nil {
			return err
		}
	}
	if qs.publisher != nil {
		if err := qs.publisher.Close(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (qs *QueueService) Publish(ctx context.Context, m interfaces.Message) error {
	if qs.publisher == nil {
		return errors.New("publisher is not defined")
	}
	return qs.publisher.Publish(ctx, m)
}
