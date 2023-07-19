package ld

import (
	"errors"
	"testing"

	"go.uber.org/zap/zapcore"
)

func TestPrefix(t *testing.T) {
	f := Prefix("session", UserAgent("1234"))
	if f.Key != "session:user-agent" {
		t.Errorf("Incorrect prefix, expected session:tsid, got %s", f.Key)
	}
}

func TestError(t *testing.T) {
	res := Error(nil)
	if res.Type != zapcore.SkipType {
		t.Errorf("Expected %v, got %v", zapcore.SkipType, res.Type)
	}
	err := errors.New("testing")
	res = Error(err)
	if res.Interface != err {
		t.Errorf("Expected %v, got %v", "testing", res.Interface)
	}
}

func TestInterfaceType(t *testing.T) {
	result := InterfaceType("str", "abc")
	if result.String != "string" {
		t.Errorf("InterfaceType: got %s, expect string", result.String)
	}
	result = InterfaceType("level", zapcore.Level(1))
	if result.String != "zapcore.Level" {
		t.Errorf("InterfaceType: got %s, expect zapcore.Level", result.String)
	}
	if result.Key != "level:type" {
		t.Errorf("InterfaceType key: got %s, expect level:type", result.Key)
	}
}

func TestURL(t *testing.T) {
	result := URL("abc")
	if result.String != "abc" {
		t.Errorf("URL: got %s, want abc", result.String)
	}
	if result.Key != "url" {
		t.Errorf("incorrect URL key: got %s", result.Key)
	}
}

func TestPort(t *testing.T) {
	result := Port(1234)
	if result.Integer != 1234 {
		t.Errorf("Port: got %d, want 1234", result.Interface)
	}
	if result.Key != "port" {
		t.Errorf("incorrect Port key: got %s", result.Key)
	}
}
func TestPortString(t *testing.T) {
	result := PortString("1234")
	if result.String != "1234" {
		t.Errorf("Port: got %s, want 1234", result.Interface)
	}
	if result.Key != "port" {
		t.Errorf("incorrect Port key: got %s", result.Key)
	}
}

func TestMethod(t *testing.T) {
	result := Method("abc")
	if result.String != "abc" {
		t.Errorf("Method: got %s, want abc", result.String)
	}
	if result.Key != "method" {
		t.Errorf("incorrect Method key: got %s", result.Key)
	}
}

func TestIP(t *testing.T) {
	result := IP("abc")
	if result.String != "abc" {
		t.Errorf("IP: got %s, want abc", result.String)
	}
	if result.Key != "ip" {
		t.Errorf("incorrect IP key: got %s", result.Key)
	}
}

func TestUserAgent(t *testing.T) {
	result := UserAgent("abc")
	if result.String != "abc" {
		t.Errorf("userAgent: got %s, want abc", result.String)
	}
	if result.Key != "user-agent" {
		t.Errorf("incorrect userAgent key: got %s", result.Key)
	}
}
