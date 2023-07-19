package logger

import "go.uber.org/zap"

// Option is a function that can be used to configure a zap.Config
type Option func(config *zap.Config)

// Debug sets the log level to debug
func Debug(config *zap.Config) { config.Level = zap.NewAtomicLevelAt(zap.DebugLevel) }

// Info sets the log level to info
func Info(config *zap.Config) { config.Level = zap.NewAtomicLevelAt(zap.InfoLevel) }

// Warn sets the log level to warn
func Warn(config *zap.Config) { config.Level = zap.NewAtomicLevelAt(zap.WarnLevel) }

// Error sets the log level to error
func Error(config *zap.Config) { config.Level = zap.NewAtomicLevelAt(zap.ErrorLevel) }

// DPanic sets the log level to dpanic
func DPanic(config *zap.Config) { config.Level = zap.NewAtomicLevelAt(zap.DPanicLevel) }

// Panic sets the log level to panic
func Panic(config *zap.Config) { config.Level = zap.NewAtomicLevelAt(zap.PanicLevel) }

// Fatal sets the log level to fatal
func Fatal(config *zap.Config) { config.Level = zap.NewAtomicLevelAt(zap.FatalLevel) }

// WithConsoleEncoding sets the encoding to console
func WithConsoleEncoding(config *zap.Config) {
	config.Encoding = "console"
	config.EncoderConfig = zap.NewDevelopmentEncoderConfig()
}

// DisableStacktrace disabled stack traces in the log output
func DisableStacktrace(config *zap.Config) {
	config.DisableStacktrace = true
}

// DisableCaller disables the caller in the log output
func DisableCaller(config *zap.Config) {
	config.DisableCaller = true
}
