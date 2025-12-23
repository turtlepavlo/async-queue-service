package middleware

import (
	"context"
	"errors"

	goqueueErrors "github.com/turtlepavlo/async-queue-service/errors"
	"github.com/turtlepavlo/async-queue-service/interfaces"
)

func mapError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, goqueueErrors.ErrInvalidMessageFormat):
		return goqueueErrors.ErrInvalidMessageFormat
	case errors.Is(err, goqueueErrors.ErrEncodingFormatNotSupported):
		return goqueueErrors.ErrEncodingFormatNotSupported
	default:
		return goqueueErrors.Error{
			Code:    goqueueErrors.UnKnownError,
			Message: err.Error(),
		}
	}
}

func PublisherDefaultErrorMapper() interfaces.PublisherMiddlewareFunc {
	return func(next interfaces.PublisherFunc) interfaces.PublisherFunc {
		return func(ctx context.Context, e interfaces.Message) (err error) {
			err = next(ctx, e)
			return mapError(err)
		}
	}
}

func InboundMessageHandlerDefaultErrorMapper() interfaces.InboundMessageHandlerMiddlewareFunc {
	return func(next interfaces.InboundMessageHandlerFunc) interfaces.InboundMessageHandlerFunc {
		return func(ctx context.Context, m interfaces.InboundMessage) (err error) {
			err = next(ctx, m)
			return mapError(err)
		}
	}
}
