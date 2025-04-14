package logger

import (
	"time"

	"go.uber.org/zap"
)

type TimedLogConfig struct {
	ErrorDuration time.Duration
	WarnDuration  time.Duration
	InfoDuration  time.Duration
	DebugDuration time.Duration
}

func DefaultTimedLogConfig() *TimedLogConfig {
	return defaultTimedLogConfig
}

func (c *TimedLogConfig) NewLog(message string, fields ...zap.Field) *TimedLog {
	return NewTimedLog(c, message, fields...)
}

type TimedLog struct {
	config   *TimedLogConfig
	message  string
	fields   []zap.Field
	duration time.Duration
	start    time.Time
	complete bool
}

func (tl *TimedLog) Complete() {
	if !tl.complete {
		tl.duration = time.Since(tl.start)
		tl.complete = true
	}
}

var defaultTimedLogConfig = &TimedLogConfig{
	ErrorDuration: time.Minute,
	WarnDuration:  30 * time.Second,
	InfoDuration:  2 * time.Second,
	DebugDuration: 500 * time.Millisecond,
}

func NewTimedLog(cnf *TimedLogConfig, message string, fields ...zap.Field) *TimedLog {
	return &TimedLog{config: cnf, message: message, fields: fields, start: time.Now()}
}

func (l *Logger) TimedLog(tl *TimedLog, fields ...zap.Field) {
	if l == nil || tl == nil {
		return
	}

	tl.Complete()

	nl := l.Clone()
	nl.zapper = nl.zapper.WithOptions(zap.AddCallerSkip(1))

	logFields := tl.fields
	logFields = append(logFields, tl.fields...)
	logFields = append(logFields, fields...)
	logFields = append(logFields, zap.Duration("duration", tl.duration))

	if tl.duration >= tl.config.ErrorDuration && tl.config.ErrorDuration > 0 {
		nl.Error(tl.message, logFields...)
	} else if tl.duration >= tl.config.WarnDuration && tl.config.WarnDuration > 0 {
		nl.Warn(tl.message, logFields...)
	} else if tl.duration >= tl.config.InfoDuration && tl.config.InfoDuration > 0 {
		nl.Info(tl.message, logFields...)
	} else if tl.duration >= tl.config.DebugDuration && tl.config.DebugDuration > 0 {
		nl.Debug(tl.message, logFields...)
	}
}
