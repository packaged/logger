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

// Clone clones the logger instance
func (l *Logger) Clone() *Logger {
	newLog := *l
	return &newLog
}

// AddCommon adds common fields to the logger
func (l *Logger) AddCommon(fields ...zap.Field) {
	l.common = append(l.common, fields...)
}

// WithCommon returns fields with common fields appended
func (l *Logger) withCommon(fields ...zap.Field) []zap.Field {
	return append(l.common, fields...)
}

// Debug logs a message at DebugLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.zapper.Debug(msg, l.withCommon(fields...)...)
}

// Info logs a message at InfoLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.zapper.Info(msg, l.withCommon(fields...)...)
}

// Warn logs a message at WarnLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.zapper.Warn(msg, l.withCommon(fields...)...)
}

// Error logs a message at ErrorLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func (l *Logger) Error(msg string, fields ...zap.Field) {
	l.zapper.Error(msg, l.withCommon(fields...)...)
}

// DPanic logs a message at DPanicLevel. The message includes any fields
// passed at the log site, as well as any fields accumulated on the logger.
//
// If the logger is in development mode, it then panics (DPanic means
// "development panic"). This is useful for catching errors that are
// recoverable, but shouldn't ever happen.
func (l *Logger) DPanic(msg string, fields ...zap.Field) {
	l.zapper.DPanic(msg, l.withCommon(fields...)...)
}

// Panic logs a message at PanicLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
//
// The logger then panics, even if logging at PanicLevel is disabled.
func (l *Logger) Panic(msg string, fields ...zap.Field) {
	l.zapper.Panic(msg, l.withCommon(fields...)...)
}

// ErrorIf logs a message at ErrorLevel if the error is not nil
func (l *Logger) ErrorIf(err error, msg string, fields ...zap.Field) {
	if err != nil {
		l.zapper.Error(msg, l.withCommon(append(fields, zap.Error(err))...)...)
	}
}

// Fatal logs a message at FatalLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
//
// The logger then calls os.Exit(1), even if logging at FatalLevel is
// disabled.
func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	l.zapper.Fatal(msg, l.withCommon(fields...)...)
}

// FatalIf logs a fatal message if the error is not nil
func (l *Logger) FatalIf(err error, msg string, fields ...zap.Field) {
	if err != nil {
		l.zapper.Fatal(msg, l.withCommon(append(fields, zap.Error(err))...)...)
	}
}

// Sync calls the underlying Core's Sync method, flushing any buffered log
// entries. Applications should take care to call Sync before exiting.
func (l *Logger) Sync() {
	_ = l.zapper.Sync()
}

// Replace replaces the global logger instance
func Replace(logger *Logger) { inst = logger }

// ObserverForTest returns an observer for the global logger instance if it exists
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
