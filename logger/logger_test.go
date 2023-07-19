package logger

import (
	"errors"
	"github.com/packaged/environment/environment"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
	"reflect"
	"testing"

	"go.uber.org/zap"
)

func TestSingleton(t *testing.T) {
	inst = nil
	i, _ := InstanceWithConfig(environment.UnitTest, zap.NewDevelopmentConfig())
	Replace(i)
	assert.NotNil(t, inst)
	reflect.DeepEqual(inst, I())

	logger := &Logger{zapper: zap.NewNop()}
	Replace(logger)
	reflect.DeepEqual(logger, I())
}

func TestGlobalClone(t *testing.T) {
	observedZapCore, observedLogs := observer.New(zap.InfoLevel)
	observedLogger := zap.New(observedZapCore)

	Replace(nil)

	assert.Panics(t, func() {
		I().Info("test")
	})

	gLog := &Logger{zapper: observedLogger}
	Replace(gLog)

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

func Test_ObserverForTest(t *testing.T) {
	Setup(environment.UnitTest)

	logObs := ObserverForTest()
	I().Debug("test message", zap.String("int1", "1"))
	if logObs.Len() != 1 {
		t.Errorf("expected 1 log, got %d", logObs.Len())
	} else if logObs.All()[0].Message != "test message" {
		t.Errorf("expected message to be 'test message', got %s", logObs.All()[0].Message)
	}

	environments := map[environment.Environment]bool{
		environment.Production:      false,
		environment.Sandbox:         false,
		environment.Development:     true,
		environment.Local:           true,
		environment.UnitTest:        true,
		environment.IntegrationTest: true,
	}

	for env, expect := range environments {
		Setup(env)
		logObs = ObserverForTest()
		if expect && logObs == nil {
			t.Errorf("expected observer for %s", env)
		}
	}
}

func TestInstanceWithConfigError(t *testing.T) {
	result, err := InstanceWithConfig(environment.UnitTest, zap.Config{})
	if err == nil {
		t.Errorf("expected error, got %v", result)
		if result.env != environment.UnitTest {
			t.Errorf("expected environment.UnitTest, got %s", result.env)
		}
	}
}

func TestSetupOptions(t *testing.T) {

	tests := []struct {
		env     environment.Environment
		options []Option
	}{
		{environment.Production, []Option{WithGoogleEncoding, DisableStacktrace, Info}},
		{environment.Sandbox, []Option{WithGoogleEncoding, DisableStacktrace, Info}},
		{environment.IntegrationTest, []Option{WithGoogleEncoding, DisableStacktrace, Debug}},
		{environment.Development, []Option{WithConsoleEncoding, Debug}},
		{environment.UnitTest, []Option{WithConsoleEncoding, Debug}},
		{environment.Local, []Option{WithConsoleEncoding, Debug}},
	}

	for _, test := range tests {
		t.Run(string(test.env)+"_options", func(t *testing.T) {
			lg, err := setup(test.env)
			if err != nil {
				t.Errorf("expected no error, got %v", err)
			} else {
				cnf := &zap.Config{}
				for _, opt := range test.options {
					opt(cnf)
				}
				setCnf := &zap.Config{}
				for _, opt := range lg.options {
					opt(setCnf)
				}

				if cnf.Level.Level() != setCnf.Level.Level() {
					t.Errorf("expected level %v, got %v", cnf.Level, setCnf.Level)
				}
				if cnf.Encoding != setCnf.Encoding {
					t.Errorf("expected encoding %v, got %v", cnf.Encoding, setCnf.Encoding)
				}
				if cnf.DisableStacktrace != setCnf.DisableStacktrace {
					t.Errorf("expected disable stacktrace %v, got %v", cnf.DisableStacktrace, setCnf.DisableStacktrace)
				}
				if cnf.DisableCaller != setCnf.DisableCaller {
					t.Errorf("expected disable caller %v, got %v", cnf.DisableCaller, setCnf.DisableCaller)
				}
			}
		})
	}
}
