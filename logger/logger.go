package logger

import (
	"github.com/packaged/environment/environment"
	"go.uber.org/zap/zaptest/observer"
	"log"

	"go.uber.org/zap"
)

var inst *Logger

// Logger is a wrapper around zap.Logger
type Logger struct {
	zapper  *zap.Logger
	options []Option
	env     environment.Environment
	common  []zap.Field
}

// I global logger instance
func I() *Logger {
	return inst
}

func setup(env environment.Environment) (zapper *Logger, err error) {
	if env.IsIntegrationTest() || BinaryDebugLogging.WithDefault("false") == "true" {
		zapper, err = InstanceWithConfig(env, zap.NewProductionConfig(), WithGoogleEncoding, DisableStacktrace, Debug)
	} else if env.IsDevOrTest() || env.IsUnitTest() {
		zapper, err = InstanceWithConfig(env, zap.NewDevelopmentConfig(), WithConsoleEncoding, Debug)
	} else {
		zapper, err = InstanceWithConfig(env, zap.NewProductionConfig(), WithGoogleEncoding, DisableStacktrace, Info)
	}
	return
}

// Setup the global logger instance
func Setup(env environment.Environment) error {
	zapper, err := setup(env)
	if err == nil {
		inst = zapper
	}
	return err
}

// InstanceWithConfig creates a new logger instance with the provided config & options applied
func InstanceWithConfig(env environment.Environment, cfg zap.Config, options ...Option) (*Logger, error) {
	// Configure with options
	for _, opt := range options {
		opt(&cfg)
	}

	zapper, err := cfg.Build()
	if err != nil {
		log.Println("Unable to create logger", err)
		return nil, err
	}
	return &Logger{env: env, zapper: zapper.WithOptions(zap.AddCallerSkip(1)), options: options}, nil
}

func (l *Logger) Clone() *Logger {
	newLog := *l
	return &newLog
}

func (l *Logger) AddCommon(fields ...zap.Field) {
	l.common = append(l.common, fields...)
}

func (l *Logger) withCommon(fields ...zap.Field) []zap.Field {
	return append(l.common, fields...)
}

func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.zapper.Debug(msg, l.withCommon(fields...)...)
}

func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.zapper.Info(msg, l.withCommon(fields...)...)
}

func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.zapper.Warn(msg, l.withCommon(fields...)...)
}

func (l *Logger) Error(msg string, fields ...zap.Field) {
	l.zapper.Error(msg, l.withCommon(fields...)...)
}

func (l *Logger) ErrorIf(err error, msg string, fields ...zap.Field) {
	if err != nil {
		l.zapper.Error(msg, l.withCommon(append(fields, zap.Error(err))...)...)
	}
}

func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	l.zapper.Fatal(msg, l.withCommon(fields...)...)
}

func (l *Logger) FatalIf(err error, msg string, fields ...zap.Field) {
	if err != nil {
		l.zapper.Fatal(msg, l.withCommon(append(fields, zap.Error(err))...)...)
	}
}

func (l *Logger) Sync() {
	_ = l.zapper.Sync()
}

func Replace(logger *Logger) { inst = logger }

func ObserverForTest() *observer.ObservedLogs {
	if inst == nil || !inst.env.IsDevOrTest() {
		return nil
	}

	observedZapCore, observedLogs := observer.New(zap.DebugLevel)
	m := Logger{env: inst.env, zapper: zap.New(observedZapCore)}
	inst.Sync()
	Replace(&m)
	return observedLogs
}
