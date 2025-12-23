package interfaces

import "context"

//go:generate mockery --name InboundMessageHandler
type InboundMessageHandler interface {
	HandleMessage(ctx context.Context, m InboundMessage) (err error)
}

type InboundMessageHandlerFunc func(ctx context.Context, m InboundMessage) (err error)

func (mhf InboundMessageHandlerFunc) HandleMessage(ctx context.Context, m InboundMessage) (err error) {
	return mhf(ctx, m)
}

type InboundMessageHandlerMiddlewareFunc func(next InboundMessageHandlerFunc) InboundMessageHandlerFunc

type InboundMessage struct {
	Message
	RetryCount int64          `json:"retryCount"`
	Metadata   map[string]any `json:"metadata"`
	// Ack confirms the message and removes it from the queue.
	Ack func(ctx context.Context) (err error) `json:"-"`
	// Nack rejects the message and requeues it for redelivery.
	Nack func(ctx context.Context) (err error) `json:"-"`
	// MoveToDeadLetterQueue rejects the message without requeueing (sends to DLQ).
	// See https://www.rabbitmq.com/docs/dlx for RabbitMQ DLX configuration.
	MoveToDeadLetterQueue func(ctx context.Context) (err error) `json:"-"`
	// RetryWithDelayFn requeues the message after a computed delay.
	RetryWithDelayFn func(ctx context.Context, delayFn DelayFn) (err error) `json:"-"`
}
