package goqueue

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/turtlepavlo/async-queue-service/internal/shared"
)

// SetupLogging configures zerolog defaults. Safe to call multiple times.
func SetupLogging() {
	shared.SetupLogging()
}

// SetupLoggingWithDefaults configures zerolog with pretty console output for development.
func SetupLoggingWithDefaults() {
	SetupLogging()
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout}).
		With().
		Caller().
		Logger()
}
