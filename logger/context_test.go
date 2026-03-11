package logger

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestFromContext_NoLogger_ReturnsGlobal(t *testing.T) {
	observedZapCore, _ := observer.New(zap.InfoLevel)
	gLog := &Logger{zapper: zap.New(observedZapCore)}
	Replace(gLog)

	got := FromContext(context.Background())
	assert.Equal(t, gLog, got)
}

func TestNewContext_RoundTrip(t *testing.T) {
	observedZapCore, _ := observer.New(zap.InfoLevel)
	l := &Logger{zapper: zap.New(observedZapCore)}

	ctx := NewContext(context.Background(), l)
	got := FromContext(ctx)
	assert.Equal(t, l, got)
}

func TestNewContext_ClonedLoggerWithCommonFields(t *testing.T) {
	observedZapCore, observedLogs := observer.New(zap.InfoLevel)
	gLog := &Logger{zapper: zap.New(observedZapCore)}
	Replace(gLog)

	cloned := gLog.Clone()
	cloned.AddCommon(zap.String("request-id", "abc-123"))

	ctx := NewContext(context.Background(), cloned)
	got := FromContext(ctx)

	got.Info("test")
	logs := observedLogs.TakeAll()
	assert.Len(t, logs, 1)
	assert.Equal(t, "test", logs[0].Message)
	assert.Len(t, logs[0].Context, 1)
	assert.Equal(t, "request-id", logs[0].Context[0].Key)
	assert.Equal(t, "abc-123", logs[0].Context[0].String)
}

func TestFromContext_NilValue_ReturnsGlobal(t *testing.T) {
	observedZapCore, _ := observer.New(zap.InfoLevel)
	gLog := &Logger{zapper: zap.New(observedZapCore)}
	Replace(gLog)

	ctx := context.WithValue(context.Background(), ctxKey{}, (*Logger)(nil))
	got := FromContext(ctx)
	assert.Equal(t, gLog, got)
}
