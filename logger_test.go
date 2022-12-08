package logger

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
	"reflect"
	"testing"

	"go.uber.org/zap"
)

func TestSingleton(t *testing.T) {
	inst = nil
	Init(DevelopmentInstance())
	assert.NotNil(t, inst)
	reflect.DeepEqual(inst, I())

	logger := &Logger{zapper: zap.NewNop()}
	Init(logger)
	reflect.DeepEqual(logger, I())
}

func TestGlobalClone(t *testing.T) {
	observedZapCore, observedLogs := observer.New(zap.InfoLevel)
	observedLogger := zap.New(observedZapCore)

	Init(nil)

	assert.Panics(t, func() {
		I().Info("test")
	})

	gLog := &Logger{zapper: observedLogger}
	Init(gLog)

	I().Info("test")
	check := observedLogs.TakeAll()
	assert.Len(t, check, 1)
	assert.Equal(t, check[0].Message, "test")

	// check global with a single common context item
	gLog.AddCommon(zap.String("test", "global1"))
	I().Info("test1")
	check = observedLogs.TakeAll()
	assert.Len(t, check, 1)
	assert.Equal(t, check[0].Message, "test1")
	assert.Len(t, check[0].Context, 1)
	assert.Equal(t, check[0].Context[0].String, "global1")

	// check clone with inherited global plus an additional common context item
	lLog := gLog.Clone()
	lLog.AddCommon(zap.String("l test", "global2"))
	lLog.Info("l test")
	check = observedLogs.TakeAll()
	assert.Len(t, check, 1)
	assert.Equal(t, check[0].Message, "l test")
	assert.Len(t, check[0].Context, 2)
	assert.Equal(t, check[0].Context[0].String, "global1")
	assert.Equal(t, check[0].Context[1].String, "global2")

	// check global still only has one (ensure that logger is correctly cloned)
	I().Info("test1")
	check = observedLogs.TakeAll()
	assert.Len(t, check, 1)
	assert.Equal(t, check[0].Message, "test1")
	assert.Len(t, check[0].Context, 1)
	assert.Equal(t, check[0].Context[0].String, "global1")
}

func TestLogging(t *testing.T) {
	fn := func(l *Logger, msg string, fields ...zap.Field) {
		l.Debug(msg, fields...)
		l.Info(msg, fields...)
		l.Warn(msg, fields...)
		l.Error(msg, fields...)
		l.Fatal(msg, fields...)
	}
	tests := []struct {
		name       string
		level      zapcore.Level
		expectLogs int
	}{
		{
			name:       "debug",
			level:      zapcore.DebugLevel,
			expectLogs: 5,
		},
		{
			name:       "info",
			level:      zapcore.InfoLevel,
			expectLogs: 4,
		},
		{
			name:       "warn",
			level:      zapcore.WarnLevel,
			expectLogs: 3,
		},
		{
			name:       "error",
			level:      zapcore.ErrorLevel,
			expectLogs: 2,
		},
		{
			name:       "fatal",
			level:      zapcore.FatalLevel,
			expectLogs: 1,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			observedZapCore, observedLogs := observer.New(test.level)
			l := &Logger{zapper: zap.New(observedZapCore, zap.WithFatalHook(zapcore.WriteThenGoexit))}
			done := make(chan interface{})
			go func() {
				defer func() {
					done <- recover()
				}()
				fn(l, "test")
				done <- true
			}()
			<-done
			logs := observedLogs.TakeAll()
			assert.Len(t, logs, test.expectLogs)
			assert.Equal(t, logs[0].Level, test.level)
		})
	}
}

func TestLogger_ErrorIf(t *testing.T) {
	observedZapCore, observedLogs := observer.New(zapcore.ErrorLevel)
	l := &Logger{zapper: zap.New(observedZapCore)}
	l.ErrorIf(nil, "test")
	logs := observedLogs.TakeAll()
	assert.Len(t, logs, 0)

	l.ErrorIf(errors.New("test"), "test")
	logs = observedLogs.TakeAll()
	assert.Len(t, logs, 1)
}

func TestLogger_FatalIf(t *testing.T) {
	observedZapCore, observedLogs := observer.New(zapcore.ErrorLevel)
	l := &Logger{zapper: zap.New(observedZapCore, zap.WithFatalHook(zapcore.WriteThenPanic))}
	l.FatalIf(nil, "test")
	logs := observedLogs.TakeAll()
	assert.Len(t, logs, 0)

	done := make(chan interface{})
	go func() {
		defer func() {
			done <- recover()
		}()
		l.FatalIf(errors.New("test"), "test")
		done <- true
	}()
	<-done
	logs = observedLogs.TakeAll()
	assert.Len(t, logs, 1)
}
