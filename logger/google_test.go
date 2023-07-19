package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"testing"
)

func TestWithGoogleEncodingOption(t *testing.T) {
	cfg := &zap.Config{}
	WithGoogleEncoding(cfg)

	if cfg.Encoding != "json" {
		t.Errorf("cfg.Encoding = %s; want json", cfg.Encoding)
	}

	encodeConfig := cfg.EncoderConfig
	if encodeConfig.TimeKey != "timestamp" {
		t.Errorf("encodeConfig.TimeKey = %s; want timestamp", encodeConfig.TimeKey)
	}
	if encodeConfig.LevelKey != "severity" {
		t.Errorf("encodeConfig.LevelKey = %s; want severity", encodeConfig.LevelKey)
	}
	if encodeConfig.NameKey != "logName" {
		t.Errorf("encodeConfig.NameKey = %s; want logName", encodeConfig.NameKey)
	}
	if encodeConfig.CallerKey != "caller" {
		t.Errorf("encodeConfig.CallerKey = %s; want caller", encodeConfig.CallerKey)
	}
	if encodeConfig.MessageKey != "textPayload" {
		t.Errorf("encodeConfig.MessageKey = %s; want textPayload", encodeConfig.MessageKey)
	}
	if encodeConfig.StacktraceKey != "trace" {
		t.Errorf("encodeConfig.StacktraceKey = %s; want trace", encodeConfig.StacktraceKey)
	}
	if encodeConfig.LineEnding != zapcore.DefaultLineEnding {
		t.Errorf("encodeConfig.LineEnding = %s; want %s", encodeConfig.LineEnding, zapcore.DefaultLineEnding)
	}
	if encodeConfig.EncodeLevel == nil {
		t.Errorf("encodeConfig.EncodeLevel = nil; want not nil")
	}
}
func TestWithGoogleEncodingOptionEncodeLevel(t *testing.T) {
	cfg := &zap.Config{}
	WithGoogleEncoding(cfg)
	tests := []struct {
		level zapcore.Level
		want  string
	}{
		{zapcore.DebugLevel, "DEBUG"},
		{zapcore.InfoLevel, "INFO"},
		{zapcore.WarnLevel, "WARNING"},
		{zapcore.ErrorLevel, "ERROR"},
		{zapcore.DPanicLevel, "CRITICAL"},
		{zapcore.PanicLevel, "ALERT"},
		{zapcore.FatalLevel, "EMERGENCY"},
	}
	for _, tt := range tests {
		enc := &primitiveArrayEncoderTest{}
		cfg.EncoderConfig.EncodeLevel(tt.level, enc)
		if enc.lastString != tt.want {
			t.Errorf("enc.lastString = %s; want %s", enc.lastString, tt.want)
		}
	}

}

type primitiveArrayEncoderTest struct {
	lastString string
}

func (p *primitiveArrayEncoderTest) AppendString(s string) {
	p.lastString = s
}
func (p *primitiveArrayEncoderTest) AppendBool(bool)             {}
func (p *primitiveArrayEncoderTest) AppendByteString([]byte)     {}
func (p *primitiveArrayEncoderTest) AppendComplex128(complex128) {}
func (p *primitiveArrayEncoderTest) AppendComplex64(complex64)   {}
func (p *primitiveArrayEncoderTest) AppendFloat64(float64)       {}
func (p *primitiveArrayEncoderTest) AppendFloat32(float32)       {}
func (p *primitiveArrayEncoderTest) AppendInt(int)               {}
func (p *primitiveArrayEncoderTest) AppendInt64(int64)           {}
func (p *primitiveArrayEncoderTest) AppendInt32(int32)           {}
func (p *primitiveArrayEncoderTest) AppendInt16(int16)           {}
func (p *primitiveArrayEncoderTest) AppendInt8(int8)             {}
func (p *primitiveArrayEncoderTest) AppendUint(uint)             {}
func (p *primitiveArrayEncoderTest) AppendUint64(uint64)         {}
func (p *primitiveArrayEncoderTest) AppendUint32(uint32)         {}
func (p *primitiveArrayEncoderTest) AppendUint16(uint16)         {}
func (p *primitiveArrayEncoderTest) AppendUint8(uint8)           {}
func (p *primitiveArrayEncoderTest) AppendUintptr(uintptr)       {}
