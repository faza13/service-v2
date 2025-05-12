package logger

import (
	"errors"
	"strings"
)

type Level int8

const (
	// LevelDebug logs are typically voluminous, and are usually disabled in
	// production.
	LevelDebug Level = iota - 1
	// LevelInfo is the default logging priority.
	LevelInfo
	// LevelWarn logs are more important than Info, but don't need individual
	// human review.
	LevelWarn
	// LevelError logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	LevelError
	// LevelPanic logs a message, then panics.
	LevelPanic
	// LevelFatal logs a message, then calls os.Exit(1).
	LevelFatal
)

func (l *Level) FromString(lvl string) error {
	switch strings.TrimSpace(strings.ToLower(lvl)) {
	case "debug":
		*l = LevelDebug
	case "info":
		*l = LevelInfo
	case "warn":
		*l = LevelWarn
	case "error":
		*l = LevelError
	case "panic":
		*l = LevelPanic
	case "fatal":
		*l = LevelFatal
	default:
		return errors.New("invalid log level")
	}
	return nil
}
