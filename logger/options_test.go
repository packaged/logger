package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"testing"
)

func TestLogLevels(t *testing.T) {
	cfg := &zap.Config{}

	tests := []struct {
		name string
		opt  Option
		want zapcore.Level
	}{
		{"Debug", Debug, zap.DebugLevel},
		{"Info", Info, zap.InfoLevel},
		{"Warn", Warn, zap.WarnLevel},
		{"Error", Error, zap.ErrorLevel},
		{"DPanic", DPanic, zap.DPanicLevel},
		{"Panic", Panic, zap.PanicLevel},
		{"Fatal", Fatal, zap.FatalLevel},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.opt(cfg)
			if cfg.Level.Level() != test.want {
				t.Errorf("cfg.Level.Level() = %s; want %s", cfg.Level.Level(), test.want)
			}
		})
	}
}

func TestWithConsoleEncoding(t *testing.T) {
	cfg := &zap.Config{}
	WithConsoleEncoding(cfg)
	if cfg.Encoding != "console" {
		t.Errorf("cfg.Encoding = %s; want console", cfg.Encoding)
	}
}

func TestDisableStacktrace(t *testing.T) {
	cfg := &zap.Config{}
	DisableStacktrace(cfg)
	if !cfg.DisableStacktrace {
		t.Errorf("cfg.DisableStacktrace = %t; want true", cfg.DisableStacktrace)
	}
}

func TestDisableCaller(t *testing.T) {
	cfg := &zap.Config{}
	DisableCaller(cfg)
	if !cfg.DisableCaller {
		t.Errorf("cfg.DisableCaller = %t; want true", cfg.DisableCaller)
	}
}
