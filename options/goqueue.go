package options

import (
	"github.com/turtlepavlo/async-queue-service/interfaces"
	"github.com/turtlepavlo/async-queue-service/internal/consumer"
	"github.com/turtlepavlo/async-queue-service/internal/publisher"
)

type GoQueueOption struct {
	NumberOfConsumer int
	Consumer         consumer.Consumer
	Publisher        publisher.Publisher
	MessageHandler   interfaces.InboundMessageHandler
}

type GoQueueOptionFunc func(opt *GoQueueOption)

func DefaultGoQueueOption() *GoQueueOption {
	return &GoQueueOption{NumberOfConsumer: 1}
}

func WithNumberOfConsumer(n int) GoQueueOptionFunc {
	return func(opt *GoQueueOption) { opt.NumberOfConsumer = n }
}

func WithConsumer(c consumer.Consumer) GoQueueOptionFunc {
	return func(opt *GoQueueOption) { opt.Consumer = c }
}

func WithPublisher(p publisher.Publisher) GoQueueOptionFunc {
	return func(opt *GoQueueOption) { opt.Publisher = p }
}

func WithMessageHandler(h interfaces.InboundMessageHandler) GoQueueOptionFunc {
	return func(opt *GoQueueOption) { opt.MessageHandler = h }
}
