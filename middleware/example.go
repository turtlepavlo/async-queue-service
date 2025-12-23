package middleware

import (
	"context"

	"github.com/rs/zerolog/log"

	"github.com/turtlepavlo/async-queue-service/interfaces"
)

// Example middlewares — copy and adapt for your own use.

func HelloWorldMiddlewareExecuteAfterInboundMessageHandler() interfaces.InboundMessageHandlerMiddlewareFunc {
	return func(next interfaces.InboundMessageHandlerFunc) interfaces.InboundMessageHandlerFunc {
		return func(ctx context.Context, m interfaces.InboundMessage) (err error) {
			err = next(ctx, m)
			if err != nil {
				log.Error().Err(err).Msg("handler error — hook in Sentry or similar here")
			}
			log.Info().Msg("hello-world-last-middleware executed")
			return err
		}
	}
}

func HelloWorldMiddlewareExecuteBeforeInboundMessageHandler() interfaces.InboundMessageHandlerMiddlewareFunc {
	return func(next interfaces.InboundMessageHandlerFunc) interfaces.InboundMessageHandlerFunc {
		return func(ctx context.Context, m interfaces.InboundMessage) (err error) {
			log.Info().Msg("hello-world-first-middleware executed")
			return next(ctx, m)
		}
	}
}

func HelloWorldMiddlewareExecuteAfterPublisher() interfaces.PublisherMiddlewareFunc {
	return func(next interfaces.PublisherFunc) interfaces.PublisherFunc {
		return func(ctx context.Context, m interfaces.Message) (err error) {
			err = next(ctx, m)
			if err != nil {
				log.Error().Err(err).Msg("publish error")
				return err
			}
			log.Info().Msg("hello-world-last-middleware executed")
			return nil
		}
	}
}

func HelloWorldMiddlewareExecuteBeforePublisher() interfaces.PublisherMiddlewareFunc {
	return func(next interfaces.PublisherFunc) interfaces.PublisherFunc {
		return func(ctx context.Context, e interfaces.Message) (err error) {
			log.Info().Msg("hello-world-first-middleware executed")
			return next(ctx, e)
		}
	}
}
