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

type Log struct {
	lg *zap.Logger
}

func (l *Log) Debug(msg string, fields ...Field) {
	l.Log(LevelDebug, msg, fields...)
}

func (l *Log) Info(msg string, fields ...Field) {
	l.Log(LevelInfo, msg, fields...)
}

func (l *Log) Warn(msg string, fields ...Field) {
	l.Log(LevelWarn, msg, fields...)
}

func (l *Log) Error(msg string, fields ...Field) {
	l.Log(LevelError, msg, fields...)
}

func (l *Log) Panic(msg string, fields ...Field) {
	l.Log(LevelPanic, msg, fields...)
}

func (l *Log) Fatal(msg string, fields ...Field) {
	l.Log(LevelFatal, msg, fields...)
}

func (l *Log) Flush() error {
	return l.lg.Sync()
}

func (l *Log) Log(lvl Level, msg string, fields ...Field) {
	zapFields := make([]zap.Field, len(fields))
	for i, field := range fields {
		zapFields[i] = zap.Any(field.Key, field.Val)
	}
	l.lg.Log(zapcore.Level(lvl), msg, zapFields...)
}

var Logger Log

func NewLogger(level Level) Log {
	lgCfg := zap.NewProductionConfig()
	if !setting.Setting.IsProduction() {
		lgCfg = zap.NewDevelopmentConfig()
	}
	lgCfg.Level.SetLevel(zapcore.Level(level))
	lg, _ := lgCfg.Build()
	Logger = Log{lg: lg}
	return Logger
}
