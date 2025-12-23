package publisher

import (
	"github.com/turtlepavlo/async-queue-service/internal/publisher"
	"github.com/turtlepavlo/async-queue-service/internal/publisher/rabbitmq"
	_ "github.com/turtlepavlo/async-queue-service/internal/shared" // Auto-setup logging
	"github.com/turtlepavlo/async-queue-service/options"
	publisherOpts "github.com/turtlepavlo/async-queue-service/options/publisher"
)

func NewPublisher(platform options.Platform, opts ...publisherOpts.PublisherOptionFunc) publisher.Publisher {
	switch platform {
	case publisherOpts.PublisherPlatformRabbitMQ:
		return rabbitmq.NewPublisher(opts...)
	case publisherOpts.PublisherPlatformGooglePubSub:
		// TODO: implement Google Pub/Sub publisher
	case publisherOpts.PublisherPlatformSNS:
		// TODO: implement SNS publisher
	default:
		panic("unknown publisher platform")
	}
	return nil
}
