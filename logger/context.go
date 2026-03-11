package logger

import "context"

type ctxKey struct{}

// NewContext returns a new context with the logger attached.
func NewContext(ctx context.Context, l *Logger) context.Context {
	return context.WithValue(ctx, ctxKey{}, l)
}

// FromContext retrieves the logger from the context.
// Returns the global logger instance if none is set on the context.
func FromContext(ctx context.Context) *Logger {
	if l, ok := ctx.Value(ctxKey{}).(*Logger); ok && l != nil {
		return l
	}
	return I()
}
