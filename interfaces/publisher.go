package interfaces

import "context"

//go:generate mockery --name PublisherHandler
type PublisherHandler interface {
	Publish(ctx context.Context, m Message) (err error)
}

type PublisherFunc func(ctx context.Context, m Message) (err error)

func (f PublisherFunc) Publish(ctx context.Context, m Message) (err error) {
	return f(ctx, m)
}

type PublisherMiddlewareFunc func(next PublisherFunc) PublisherFunc
