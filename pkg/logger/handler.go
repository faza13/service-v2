package logger

import (
	"context"

	"go.uber.org/zap/zapcore"
)

// zapHandler provides a wrapper around the zap core to capture which
// log level is being logged for event handling.
type zapHandler struct {
	core   zapcore.Core
	events Events
}

func newZapHandler(core zapcore.Core, events Events) *zapHandler {
	return &zapHandler{
		core:   core,
		events: events,
	}
}

// Enabled reports whether the handler handles entries at the given level.
// The handler ignores entries whose level is lower.
func (h *zapHandler) Enabled(level zapcore.Level) bool {
	return h.core.Enabled(level)
}

// With returns a new Core that includes the given fields in each output entry.
func (h *zapHandler) With(fields []zapcore.Field) zapcore.Core {
	return &zapHandler{core: h.core.With(fields), events: h.events}
}

// Check determines whether the supplied Entry should be logged.
func (h *zapHandler) Check(entry zapcore.Entry, checked *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if h.Enabled(entry.Level) {
		return checked.AddCore(entry, h)
	}
	return checked
}

// Write serializes the Entry and any Fields supplied at the log site and
// writes them to the underlying Core.
func (h *zapHandler) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	// Process events based on log level
	switch entry.Level {
	case zapcore.DebugLevel:
		if h.events.Debug != nil {
			h.events.Debug(context.Background(), toRecord(entry, fields))
		}

	case zapcore.ErrorLevel:
		if h.events.Error != nil {
			h.events.Error(context.Background(), toRecord(entry, fields))
		}

	case zapcore.WarnLevel:
		if h.events.Warn != nil {
			h.events.Warn(context.Background(), toRecord(entry, fields))
		}

	case zapcore.InfoLevel:
		if h.events.Info != nil {
			h.events.Info(context.Background(), toRecord(entry, fields))
		}
	}

	// Pass to the underlying core
	return h.core.Write(entry, fields)
}

// Sync flushes buffered logs (if any).
func (h *zapHandler) Sync() error {
	return h.core.Sync()
}
