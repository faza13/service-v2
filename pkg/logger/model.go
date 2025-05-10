package logger

import (
	"context"
	"time"

	"go.uber.org/zap/zapcore"
)

// Level represents different logging levels.
type Level zapcore.Level

// A set of possible logging levels.
const (
	LevelDebug = Level(zapcore.DebugLevel)
	LevelInfo  = Level(zapcore.InfoLevel)
	LevelWarn  = Level(zapcore.WarnLevel)
	LevelError = Level(zapcore.ErrorLevel)
)

// Record represents the data that is being logged.
type Record struct {
	Time       time.Time
	Message    string
	Level      Level
	Attributes map[string]any
}

func toRecord(entry zapcore.Entry, fields []zapcore.Field) Record {
	atts := make(map[string]any, len(fields))

	for _, field := range fields {
		atts[field.Key] = fieldToAny(field)
	}

	return Record{
		Time:       entry.Time,
		Message:    entry.Message,
		Level:      Level(entry.Level),
		Attributes: atts,
	}
}

// fieldToAny converts a zapcore.Field to its underlying value
func fieldToAny(field zapcore.Field) any {
	switch field.Type {
	case zapcore.StringType:
		return field.String
	case zapcore.BoolType:
		return field.Integer == 1
	case zapcore.Int64Type, zapcore.Int32Type, zapcore.Int16Type, zapcore.Int8Type:
		return field.Integer
	case zapcore.Uint64Type, zapcore.Uint32Type, zapcore.Uint16Type, zapcore.Uint8Type:
		return uint64(field.Integer)
	case zapcore.Float64Type:
		return float64(field.Integer)
	case zapcore.Float32Type:
		return float32(field.Integer)
	case zapcore.DurationType:
		return time.Duration(field.Integer)
	case zapcore.TimeType:
		if field.Interface != nil {
			return field.Interface.(time.Time)
		}
		return time.Unix(0, field.Integer)
	case zapcore.ReflectType:
		return field.Interface
	default:
		return field.Interface
	}
}

// EventFn is a function to be executed when configured against a log level.
type EventFn func(ctx context.Context, r Record)

// Events contains an assignment of an event function to a log level.
type Events struct {
	Debug EventFn
	Info  EventFn
	Warn  EventFn
	Error EventFn
}
