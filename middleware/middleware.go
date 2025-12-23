package middleware

import (
	"github.com/turtlepavlo/async-queue-service/interfaces"
)

func ApplyHandlerMiddleware(h interfaces.InboundMessageHandlerFunc,
	middlewares ...interfaces.InboundMessageHandlerMiddlewareFunc) interfaces.InboundMessageHandlerFunc {
	for _, middleware := range middlewares {
		h = middleware(h)
	}
	return h
}

func ApplyPublisherMiddleware(p interfaces.PublisherFunc,
	middlewares ...interfaces.PublisherMiddlewareFunc) interfaces.PublisherFunc {
	for _, middleware := range middlewares {
		p = middleware(p)
	}
	return p
}
