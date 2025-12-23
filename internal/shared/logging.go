package shared

import (
	"sync"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

var loggingSetupOnce sync.Once

// SetupLogging configures zerolog defaults (Unix timestamps + stack traces). Runs once.
func SetupLogging() {
	loggingSetupOnce.Do(func() {
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	})
}

//nolint:gochecknoinits // auto-setup logging when any goqueue package is imported
func init() {
	SetupLogging()
}
