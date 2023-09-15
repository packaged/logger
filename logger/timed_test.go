package logger

import (
	"github.com/packaged/environment/environment"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"testing"
	"time"
)

func TestTimedLog(t *testing.T) {

	Setup(environment.UnitTest)
	logObs := ObserverForTest()
	I().TimedLog(nil)
	if logObs.Len() > 0 {
		t.Errorf("expected no logs, got %d", logObs.Len())
	}

	tl := DefaultTimedLogConfig().NewLog("Doing something")
	time.Sleep(time.Second)
	I().TimedLog(tl)
	if logObs.Len() < 1 {
		t.Errorf("expected 1 log, got %d", logObs.Len())
	}

	logs := logObs.All()
	timedLog := logs[0]
	for _, f := range timedLog.Context {
		if f.Key == "duration" {
			if f.Integer < time.Second.Nanoseconds() {
				t.Errorf("expected at least 1 second, got %dns", f.Integer)
			}
		}
	}
}

func TestTimedLogBreaks(t *testing.T) {

	cnf := &TimedLogConfig{
		ErrorDuration: time.Millisecond * 30,
		WarnDuration:  time.Millisecond * 20,
		InfoDuration:  time.Millisecond * 10,
		DebugDuration: time.Millisecond,
	}

	tests := []struct {
		name     string
		duration time.Duration
		level    zapcore.Level
	}{
		{"error", cnf.ErrorDuration, zapcore.ErrorLevel},
		{"warn", cnf.WarnDuration, zapcore.WarnLevel},
		{"info", cnf.InfoDuration, zapcore.InfoLevel},
		{"debug", cnf.DebugDuration, zapcore.DebugLevel},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testInst, _ := InstanceWithConfig(environment.UnitTest, zap.NewDevelopmentConfig(), WithConsoleEncoding, Debug)
			Replace(testInst)
			logObs := ObserverForTest()
			timedL := cnf.NewLog("abc")
			time.Sleep(test.duration + time.Millisecond)
			I().TimedLog(timedL)
			if logObs.Len() != 1 {
				t.Errorf("expected 1 log, got %d", logObs.Len())
			}
			logs := logObs.All()
			if logs[0].Level != test.level {
				t.Errorf("expected level %s, got %s", test.level, logs[0].Level)
			}
		})
	}
}
