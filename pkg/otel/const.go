package otel

type ctxKey int

const (
	tracerKey ctxKey = iota + 1
	traceIDKey
)
