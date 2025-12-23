package consumer

import (
	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/turtlepavlo/async-queue-service/interfaces"
	"github.com/turtlepavlo/async-queue-service/options"
)

const (
	DefaultMaxRetryFailedMessage = 3
	DefaultBatchMessageSize      = 1
)

type ConsumerOption struct {
	BatchMessageSize       int
	QueueName              string
	Middlewares            []interfaces.InboundMessageHandlerMiddlewareFunc
	MaxRetryFailedMessage  int64
	ConsumerID             string
	RabbitMQConsumerConfig *RabbitMQConsumerConfig
}

type ConsumerOptionFunc func(opt *ConsumerOption)

func WithBatchMessageSize(n int) ConsumerOptionFunc {
	return func(opt *ConsumerOption) { opt.BatchMessageSize = n }
}

func WithQueueName(name string) ConsumerOptionFunc {
	return func(opt *ConsumerOption) { opt.QueueName = name }
}

func WithMiddlewares(middlewares ...interfaces.InboundMessageHandlerMiddlewareFunc) ConsumerOptionFunc {
	return func(opt *ConsumerOption) { opt.Middlewares = middlewares }
}

func WithMaxRetryFailedMessage(n int64) ConsumerOptionFunc {
	return func(opt *ConsumerOption) { opt.MaxRetryFailedMessage = n }
}

func WithConsumerID(id string) ConsumerOptionFunc {
	return func(opt *ConsumerOption) { opt.ConsumerID = id }
}

func WithRabbitMQConsumerConfig(rabbitMQOption *RabbitMQConsumerConfig) ConsumerOptionFunc {
	return func(opt *ConsumerOption) { opt.RabbitMQConsumerConfig = rabbitMQOption }
}

var DefaultConsumerOption = func() *ConsumerOption {
	return &ConsumerOption{
		Middlewares:           []interfaces.InboundMessageHandlerMiddlewareFunc{},
		BatchMessageSize:      DefaultBatchMessageSize,
		MaxRetryFailedMessage: DefaultMaxRetryFailedMessage,
	}
}

type RabbitMQConsumerConfig struct {
	ConsumerChannel    *amqp.Channel
	ReQueueChannel     *amqp.Channel
	QueueDeclareConfig *RabbitMQQueueDeclareConfig
	QueueBindConfig    *RabbitMQQueueBindConfig
}

// RabbitMQQueueDeclareConfig mirrors the parameters of amqp.Channel.QueueDeclare.
// See https://www.rabbitmq.com/amqp-0-9-1-reference.html#queue.declare for full semantics.
type RabbitMQQueueDeclareConfig struct {
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	Args       amqp.Table
}

type RabbitMQQueueBindConfig struct {
	RoutingKeys  []string
	ExchangeName string
	NoWait       bool
	Args         amqp.Table
}

const (
	ConsumerPlatformRabbitMQ     = options.PlatformRabbitMQ
	ConsumerPlatformGooglePubSub = options.PlatformGooglePubSub
	ConsumerPlatformSQS          = options.PlatformSQS
)

// RabbitMQConfigWithDefaultTopicFanOutPattern builds a consumer config for a durable topic exchange.
// Routing keys support wildcard patterns: "payments.#" matches all sub-keys.
// See https://www.rabbitmq.com/tutorials/tutorial-five-go.html for pattern syntax.
func RabbitMQConfigWithDefaultTopicFanOutPattern(consumerChannel, requeueChannel *amqp.Channel,
	exchangeName string, routingKeys []string) *RabbitMQConsumerConfig {
	return &RabbitMQConsumerConfig{
		ConsumerChannel: consumerChannel,
		ReQueueChannel:  requeueChannel,
		QueueDeclareConfig: &RabbitMQQueueDeclareConfig{
			Durable: true,
		},
		QueueBindConfig: &RabbitMQQueueBindConfig{
			RoutingKeys:  routingKeys,
			ExchangeName: exchangeName,
		},
	}
}
