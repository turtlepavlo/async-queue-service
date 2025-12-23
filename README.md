# async-queue-service

A Go library that provides a unified abstraction over message queue platforms. Currently supports **RabbitMQ**, with Google Pub/Sub and AWS SNS/SQS planned.

## Installation

```bash
go get github.com/turtlepavlo/async-queue-service
```

Requires Go 1.24+.

## Core concepts

| Concept | Description |
|---|---|
| `QueueService` | Top-level entry point. Wires consumer, publisher and handler together. |
| `Message` | Envelope sent and received. `Action` = routing key, `Topic` = exchange. |
| `InboundMessage` | Received message with `Ack`, `Nack`, `MoveToDeadLetterQueue`, `RetryWithDelayFn`. |
| `Middleware` | Functions that wrap the handler or publisher — applied in order. |
| `Encoding` | Pluggable encode/decode per content-type. JSON is registered by default. |

## Quick start

### Publisher

```go
package main

import (
    "context"

    amqp "github.com/rabbitmq/amqp091-go"

    goqueue "github.com/turtlepavlo/async-queue-service"
    "github.com/turtlepavlo/async-queue-service/interfaces"
    "github.com/turtlepavlo/async-queue-service/options"
    "github.com/turtlepavlo/async-queue-service/publisher"
    publisherOpts "github.com/turtlepavlo/async-queue-service/options/publisher"
)

func main() {
    conn, _ := amqp.Dial("amqp://guest:guest@localhost:5672/")

    pub := publisher.NewPublisher(
        publisherOpts.PublisherPlatformRabbitMQ,
        publisherOpts.WithPublisherID("my-service"),
        publisherOpts.WithRabbitMQPublisherConfig(&publisherOpts.RabbitMQPublisherConfig{
            Conn:                     conn,
            PublisherChannelPoolSize: 5,
        }),
    )

    svc := goqueue.NewQueueService(options.WithPublisher(pub))

    _ = svc.Publish(context.Background(), interfaces.Message{
        Topic:  "my-exchange",
        Action: "order.created",
        Data:   map[string]any{"order_id": 42},
    })
}
```

### Consumer

```go
package main

import (
    "context"
    "fmt"

    amqp "github.com/rabbitmq/amqp091-go"

    goqueue "github.com/turtlepavlo/async-queue-service"
    "github.com/turtlepavlo/async-queue-service/consumer"
    "github.com/turtlepavlo/async-queue-service/interfaces"
    "github.com/turtlepavlo/async-queue-service/options"
    consumerOpts "github.com/turtlepavlo/async-queue-service/options/consumer"
)

func main() {
    conn, _ := amqp.Dial("amqp://guest:guest@localhost:5672/")
    ch, _ := conn.Channel()
    requeueCh, _ := conn.Channel()

    cons := consumer.NewConsumer(
        consumerOpts.ConsumerPlatformRabbitMQ,
        consumerOpts.WithRabbitMQConsumerConfig(
            consumerOpts.RabbitMQConfigWithDefaultTopicFanOutPattern(
                ch, requeueCh,
                "my-exchange",
                []string{"order.#"},
            ),
        ),
        consumerOpts.WithQueueName("orders-queue"),
        consumerOpts.WithBatchMessageSize(10),
        consumerOpts.WithMaxRetryFailedMessage(3),
    )

    svc := goqueue.NewQueueService(
        options.WithConsumer(cons),
        options.WithMessageHandler(handler()),
        options.WithNumberOfConsumer(2),
    )

    // blocks until ctx is cancelled
    _ = svc.Start(context.Background())
}

func handler() interfaces.InboundMessageHandlerFunc {
    return func(ctx context.Context, m interfaces.InboundMessage) error {
        fmt.Printf("received: action=%s data=%v\n", m.Action, m.Data)
        return m.Ack(ctx)
    }
}
```

## Retry mechanism (RabbitMQ)

Call `m.RetryWithDelayFn(ctx, delayFn)` inside your handler to requeue with a delay. The library uses a per-retry-count TTL queue backed by a dead-letter exchange so that messages are delivered back to the original queue only after the delay expires.

Three built-in delay functions:

| Function | Delay formula |
|---|---|
| `interfaces.ExponentialBackoffDelayFn` | `2^(retries-1)` seconds |
| `interfaces.LinearDelayFn` | `retries` seconds |
| `interfaces.NoDelayFn` | 0 seconds |

```go
func handler() interfaces.InboundMessageHandlerFunc {
    return func(ctx context.Context, m interfaces.InboundMessage) error {
        if err := process(m); err != nil {
            // retry with exponential backoff; moves to DLQ after MaxRetryFailedMessage
            return m.RetryWithDelayFn(ctx, interfaces.ExponentialBackoffDelayFn)
        }
        return m.Ack(ctx)
    }
}
```

When `RetryCount` exceeds `MaxRetryFailedMessage` (default: 3), the consumer automatically nacks the message without requeue.

## Middleware

Middlewares wrap the handler or publisher and are applied in declaration order.

```go
// Handler middleware
cons := consumer.NewConsumer(
    consumerOpts.ConsumerPlatformRabbitMQ,
    consumerOpts.WithMiddlewares(
        middleware.HelloWorldMiddlewareExecuteBeforeInboundMessageHandler(),
        middleware.HelloWorldMiddlewareExecuteAfterInboundMessageHandler(),
    ),
    // ...
)

// Publisher middleware
pub := publisher.NewPublisher(
    publisherOpts.PublisherPlatformRabbitMQ,
    publisherOpts.WithMiddlewares(
        middleware.HelloWorldMiddlewareExecuteBeforePublisher(),
        middleware.HelloWorldMiddlewareExecuteAfterPublisher(),
    ),
    // ...
)
```

Use `middleware.ApplyHandlerMiddleware` / `middleware.ApplyPublisherMiddleware` to chain middlewares manually.

Built-in middleware:
- `PublisherDefaultErrorMapper()` — maps raw errors to typed `goqueue/errors` values.
- `InboundMessageHandlerDefaultErrorMapper()` — same for inbound handlers.

## Custom encoding

JSON is registered by default. Register additional encodings before publishing:

```go
goqueue.AddGoQueueEncoding(headerVal.ContentTypeXML, &goqueue.Encoding{
    ContentType: headerVal.ContentTypeXML,
    Encode:      myXMLEncoder,
    Decode:      myXMLDecoder,
})
```

## Message headers

All messages carry a set of `goqueue-` prefixed headers automatically. **Do not use the `goqueue-` prefix** in your own headers — they will be overwritten.

| Header key | Description |
|---|---|
| `goqueue-app-id` | Publisher ID |
| `goqueue-message-id` | Auto-generated UUID |
| `goqueue-published-timestamp` | RFC3339 timestamp |
| `goqueue-retry-count` | Number of retries |
| `goqueue-content-type` | Content-type of the payload |
| `goqueue-queue-service-agent` | Platform identifier (e.g. `goqueue/rabbitmq`) |
| `goqueue-schema-version` | Schema version (`1.0.0`) |

## Configuration reference

### QueueService options

| Option | Default | Description |
|---|---|---|
| `WithNumberOfConsumer(n)` | 1 | Number of concurrent consumer goroutines |
| `WithConsumer(c)` | — | Consumer implementation |
| `WithPublisher(p)` | — | Publisher implementation |
| `WithMessageHandler(h)` | — | Message handler |

### Consumer options (RabbitMQ)

| Option | Default | Description |
|---|---|---|
| `WithQueueName(name)` | — | Queue to consume from |
| `WithBatchMessageSize(n)` | 1 | Prefetch count (messages in-flight per consumer) |
| `WithMaxRetryFailedMessage(n)` | 3 | Max retries before dropping to DLQ |
| `WithConsumerID(id)` | auto UUID | Consumer tag |
| `WithMiddlewares(...)` | — | Handler middlewares |
| `WithRabbitMQConsumerConfig(cfg)` | — | Raw RabbitMQ channel config |

### Publisher options (RabbitMQ)

| Option | Default | Description |
|---|---|---|
| `WithPublisherID(id)` | auto UUID | Identifies this publisher in message headers |
| `WithMiddlewares(...)` | — | Publisher middlewares |
| `WithRabbitMQPublisherConfig(cfg)` | — | Connection and channel pool config |

## Logging

The library uses [zerolog](https://github.com/rs/zerolog). Logging is auto-configured when any package is imported. For development, call:

```go
goqueue.SetupLoggingWithDefaults() // pretty console output with caller info
```

## Running integration tests

Integration tests require a running RabbitMQ instance:

```bash
docker compose -f test.compose.yaml up -d
RABBITMQ_TEST_URL=amqp://test:test@localhost:5672/test go test ./...
```

Short-mode skips integration tests:

```bash
go test -short ./...
```
