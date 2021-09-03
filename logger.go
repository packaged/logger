package logger

import (
	"go.uber.org/zap"
	"log"
)

func Instance(options ...Option) *zap.Logger {
	return InstanceWithConfig(zap.NewProductionConfig(), options...)
}
func DevelopmentInstance(options ...Option) *zap.Logger {
	return InstanceWithConfig(zap.NewDevelopmentConfig(), options...)
}

func InstanceWithConfig(cfg zap.Config, options ...Option) *zap.Logger {
	// Configure with options
	for _, opt := range options {
		opt(&cfg)
	}

	logger, err := cfg.Build()
	if err != nil {
		log.Println("Unable to create logger", err)
	}
	return logger

}
