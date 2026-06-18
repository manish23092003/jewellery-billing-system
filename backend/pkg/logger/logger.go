package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Init configures the global zerolog logger with a human-friendly console writer.
// Call this once at application startup before any logging occurs.
func Init() {
	zerolog.TimeFieldFormat = time.RFC3339

	log.Logger = zerolog.New(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "2006-01-02 15:04:05",
	}).With().Timestamp().Caller().Logger()
}

// Convenience wrappers — re-export zerolog levels so callers don't need
// to import zerolog directly.

func Info() *zerolog.Event  { return log.Info() }
func Error() *zerolog.Event { return log.Error() }
func Fatal() *zerolog.Event { return log.Fatal() }
func Debug() *zerolog.Event { return log.Debug() }
func Warn() *zerolog.Event  { return log.Warn() }
