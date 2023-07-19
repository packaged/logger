package logger

import "go.uber.org/zap"

type Option func(config *zap.Config)

func Debug(config *zap.Config)  { config.Level = zap.NewAtomicLevelAt(zap.DebugLevel) }
func Info(config *zap.Config)   { config.Level = zap.NewAtomicLevelAt(zap.InfoLevel) }
func Warn(config *zap.Config)   { config.Level = zap.NewAtomicLevelAt(zap.WarnLevel) }
func Error(config *zap.Config)  { config.Level = zap.NewAtomicLevelAt(zap.ErrorLevel) }
func DPanic(config *zap.Config) { config.Level = zap.NewAtomicLevelAt(zap.DPanicLevel) }
func Panic(config *zap.Config)  { config.Level = zap.NewAtomicLevelAt(zap.PanicLevel) }
func Fatal(config *zap.Config)  { config.Level = zap.NewAtomicLevelAt(zap.FatalLevel) }

func WithConsoleEncoding(config *zap.Config) {
	config.Encoding = "console"
	config.EncoderConfig = zap.NewDevelopmentEncoderConfig()
}

func DisableStacktrace(config *zap.Config) {
	config.DisableStacktrace = true
}

func DisableCaller(config *zap.Config) {
	config.DisableCaller = true
}
