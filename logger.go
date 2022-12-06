package logger

import (
	"log"

	"go.uber.org/zap"
)

var inst *Logger

type Logger struct {
	zapper *zap.Logger
	common []zap.Field
}

func I() *Logger {
	return inst
}

func Instance(options ...Option) *Logger {
	return instanceWithConfig(zap.NewProductionConfig(), options...)
}

func DevelopmentInstance(options ...Option) *Logger {
	return instanceWithConfig(zap.NewDevelopmentConfig(), options...)
}

func instanceWithConfig(cfg zap.Config, options ...Option) *Logger {
	// Configure with options
	for _, opt := range options {
		opt(&cfg)
	}

	zapper, err := cfg.Build()
	if err != nil {
		log.Println("Unable to create logger", err)
	}
	return &Logger{zapper: zapper.WithOptions(zap.AddCallerSkip(1))}
}

func Init(logger *Logger) {
	inst = logger
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
