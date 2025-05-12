package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"service/pkg/setting"
)

type Field struct {
	Key string
	Val interface{}
}

func F(key string, val interface{}) Field {
	return Field{Key: key, Val: val}
}

type Logger struct {
	lg *zap.Logger
}

func (l Logger) Debug(msg string, fields ...Field) {
	l.Log(LevelDebug, msg, fields...)
}

func (l Logger) Info(msg string, fields ...Field) {
	l.Log(LevelInfo, msg, fields...)
}

func (l Logger) Warn(msg string, fields ...Field) {
	l.Log(LevelWarn, msg, fields...)
}

func (l Logger) Error(msg string, fields ...Field) {
	l.Log(LevelError, msg, fields...)
}

func (l Logger) Panic(msg string, fields ...Field) {
	l.Log(LevelPanic, msg, fields...)
}

func (l Logger) Fatal(msg string, fields ...Field) {
	l.Log(LevelFatal, msg, fields...)
}

func (l Logger) Flush() error {
	return l.lg.Sync()
}

func (l Logger) Log(lvl Level, msg string, fields ...Field) {
	zapFields := make([]zap.Field, len(fields))
	for i, field := range fields {
		zapFields[i] = zap.Any(field.Key, field.Val)
	}
	l.lg.Log(zapcore.Level(lvl), msg, zapFields...)
}

func NewLogger(level Level) Logger {
	lgCfg := zap.NewProductionConfig()
	if !setting.Setting.IsProduction() {
		lgCfg = zap.NewDevelopmentConfig()
	}
	lgCfg.Level.SetLevel(zapcore.Level(level))
	lg, _ := lgCfg.Build()

	return Logger{lg: lg}
}
