package goqueue

import (
	"context"
	"encoding/json"
	"sync"

	headerVal "github.com/turtlepavlo/async-queue-service/headers/value"
	"github.com/turtlepavlo/async-queue-service/interfaces"
)

type EncoderFn func(ctx context.Context, m interfaces.Message) (data []byte, err error)
type DecoderFn func(ctx context.Context, data []byte) (m interfaces.Message, err error)

var (
	JSONEncoder EncoderFn = func(_ context.Context, m interfaces.Message) ([]byte, error) {
		return json.Marshal(m)
	}
	JSONDecoder DecoderFn = func(_ context.Context, data []byte) (m interfaces.Message, err error) {
		err = json.Unmarshal(data, &m)
		return
	}

	DefaultEncoder EncoderFn = JSONEncoder
	DefaultDecoder DecoderFn = JSONDecoder
)

var goQueueEncodingMap = sync.Map{}

// AddGoQueueEncoding registers a custom encoding for a content type.
// Must be called before any message with that content type is published.
func AddGoQueueEncoding(contentType headerVal.ContentType, encoding *Encoding) {
	goQueueEncodingMap.Store(contentType, encoding)
}

func GetGoQueueEncoding(contentType headerVal.ContentType) (res *Encoding, ok bool) {
	if encoding, ok := goQueueEncodingMap.Load(contentType); ok {
		return encoding.(*Encoding), ok
	}
	return nil, false
}

// Encoding pairs an encoder and decoder for a given content type.
type Encoding struct {
	ContentType headerVal.ContentType
	Encode      EncoderFn
	Decode      DecoderFn
}

var (
	JSONEncoding = &Encoding{
		ContentType: headerVal.ContentTypeJSON,
		Encode:      JSONEncoder,
		Decode:      JSONDecoder,
	}
	DefaultEncoding = JSONEncoding
)

//nolint:gochecknoinits // auto-registers JSON encoding on package import
func init() {
	AddGoQueueEncoding(JSONEncoding.ContentType, JSONEncoding)
}
