package consumer

import (
	"github.com/turtlepavlo/async-queue-service/internal/consumer"
	"github.com/turtlepavlo/async-queue-service/internal/consumer/rabbitmq"
	_ "github.com/turtlepavlo/async-queue-service/internal/shared" // Auto-setup logging
	"github.com/turtlepavlo/async-queue-service/options"
	consumerOpts "github.com/turtlepavlo/async-queue-service/options/consumer"
)

func NewConsumer(platform options.Platform, opts ...consumerOpts.ConsumerOptionFunc) consumer.Consumer {
	switch platform {
	case consumerOpts.ConsumerPlatformRabbitMQ:
		return rabbitmq.NewConsumer(opts...)
	case consumerOpts.ConsumerPlatformGooglePubSub:
		// TODO: implement Google Pub/Sub consumer
	case consumerOpts.ConsumerPlatformSQS:
		// TODO: implement SQS consumer
	default:
		panic("unknown consumer platform")
	}
	return nil
}
