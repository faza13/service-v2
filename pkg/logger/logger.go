// Package logger provides support for initializing the log system.
package logger

import (
	"context"
	"io"
	"log"
	"runtime"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// TraceIDFn represents a function that can return the trace id from
// the specified context.
type TraceIDFn func(ctx context.Context) string

// Logger represents a logger for logging information.
type Logger struct {
	logger    *zap.Logger
	core      zapcore.Core
	traceIDFn TraceIDFn
}

// New constructs a new log for application use.
func New(w io.Writer, minLevel Level, serviceName string, traceIDFn TraceIDFn) *Logger {
	return new(w, minLevel, serviceName, traceIDFn, Events{})
}

// NewWithEvents constructs a new log for application use with events.
func NewWithEvents(w io.Writer, minLevel Level, serviceName string, traceIDFn TraceIDFn, events Events) *Logger {
	return new(w, minLevel, serviceName, traceIDFn, events)
}

// NewWithCore returns a new log for application use with the underlying
// core.
func NewWithCore(core zapcore.Core) *Logger {
	logger := zap.New(core)
	return &Logger{
		logger: logger,
		core:   core,
	}
}

// NewStdLogger returns a standard library Logger that wraps the zap Logger.
func NewStdLogger(logger *Logger, level Level) *log.Logger {
	return zap.NewStdLog(logger.logger)
}

// Debug logs at LevelDebug with the given context.
func (log *Logger) Debug(ctx context.Context, msg string, args ...any) {
	log.write(ctx, LevelDebug, 3, msg, args...)
}

// Debugc logs the information at the specified call stack position.
func (log *Logger) Debugc(ctx context.Context, caller int, msg string, args ...any) {
	log.write(ctx, LevelDebug, caller, msg, args...)
}

// Info logs at LevelInfo with the given context.
func (log *Logger) Info(ctx context.Context, msg string, args ...any) {
	log.write(ctx, LevelInfo, 3, msg, args...)
}

// Infoc logs the information at the specified call stack position.
func (log *Logger) Infoc(ctx context.Context, caller int, msg string, args ...any) {
	log.write(ctx, LevelInfo, caller, msg, args...)
}

// Warn logs at LevelWarn with the given context.
func (log *Logger) Warn(ctx context.Context, msg string, args ...any) {
	log.write(ctx, LevelWarn, 3, msg, args...)
}

// Warnc logs the information at the specified call stack position.
func (log *Logger) Warnc(ctx context.Context, caller int, msg string, args ...any) {
	log.write(ctx, LevelWarn, caller, msg, args...)
}

// Error logs at LevelError with the given context.
func (log *Logger) Error(ctx context.Context, msg string, args ...any) {
	log.write(ctx, LevelError, 3, msg, args...)
}

// Errorc logs the information at the specified call stack position.
func (log *Logger) Errorc(ctx context.Context, caller int, msg string, args ...any) {
	log.write(ctx, LevelError, caller, msg, args...)
}

func (log *Logger) write(ctx context.Context, level Level, caller int, msg string, args ...any) {
	zapLevel := zapcore.Level(level)

	if !log.core.Enabled(zapLevel) {
		return
	}

	// Convert args to zap fields
	fields := argsToFields(args...)

	// Add trace ID if available
	if log.traceIDFn != nil {
		fields = append(fields, zap.String("trace_id", log.traceIDFn(ctx)))
	}

	// Get caller information
	var pcs [1]uintptr
	runtime.Callers(caller, pcs[:])
	pc := pcs[0]

	// Create entry
	entry := zapcore.Entry{
		Level:   zapLevel,
		Time:    time.Now(),
		Message: msg,
		Caller:  zapcore.EntryCaller{Defined: true, PC: pc},
		Stack:   "",
	}

	// Check and write
	ce := log.core.Check(entry, nil)
	if ce != nil {
		ce.Write(fields...)
	}
}

func argsToFields(args ...any) []zapcore.Field {
	fields := make([]zapcore.Field, 0, len(args))

	for i := 0; i < len(args); i += 2 {
		// Make sure we have both key and value
		if i+1 >= len(args) {
			break
		}

		// Key must be a string
		key, ok := args[i].(string)
		if !ok {
			continue
		}

		// Convert value to appropriate field type
		value := args[i+1]
		fields = append(fields, anyToField(key, value))
	}

	return fields
}

func anyToField(key string, value any) zapcore.Field {
	switch v := value.(type) {
	case string:
		return zap.String(key, v)
	case int:
		return zap.Int(key, v)
	case int64:
		return zap.Int64(key, v)
	case uint:
		return zap.Uint(key, v)
	case uint64:
		return zap.Uint64(key, v)
	case float64:
		return zap.Float64(key, v)
	case float32:
		return zap.Float32(key, v)
	case bool:
		return zap.Bool(key, v)
	case time.Time:
		return zap.Time(key, v)
	case time.Duration:
		return zap.Duration(key, v)
	case error:
		return zap.Error(v)
	default:
		return zap.Any(key, v)
	}
}

func new(w io.Writer, minLevel Level, serviceName string, traceIDFn TraceIDFn, events Events) *Logger {
	// Create encoder config
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Create JSON encoder
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	// Create writer syncer
	writeSyncer := zapcore.AddSync(w)

	// Create core
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.Level(minLevel))

	// If events are to be processed, wrap the core with our custom handler
	if events.Debug != nil || events.Info != nil || events.Warn != nil || events.Error != nil {
		core = newZapHandler(core, events)
	}

	// Add service name as a field to all logs
	fields := []zapcore.Field{
		zap.String("service", serviceName),
	}
	core = core.With(fields)

	// Create logger
	logger := zap.New(core, zap.AddCaller())

	return &Logger{
		logger:    logger,
		core:      core,
		traceIDFn: traceIDFn,
	}
}
