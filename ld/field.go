package ld

// Log Data

import (
	"reflect"

	"go.uber.org/zap"
)

type Field zap.Field

func Prefix(prefix string, field zap.Field) zap.Field {
	field.Key = prefix + ":" + field.Key
	return field
}

func Error(err error) zap.Field {
	if err == nil {
		return zap.Skip()
	}
	return zap.Error(err)
}

func InterfaceType(key string, iface interface{}) zap.Field {
	return zap.String(key+":type", reflect.TypeOf(iface).String())
}

func IP(ip string) zap.Field        { return zap.String("ip", ip) }
func UserAgent(ua string) zap.Field { return zap.String("user-agent", ua) }

func URL(input string) zap.Field                { return zap.String("url", input) }
func PortString(input string) zap.Field         { return zap.String("port", input) }
func Port(input int) zap.Field                  { return zap.Int("port", input) }
func Method(input string) zap.Field             { return zap.String("method", input) }
func TrustedString(key, input string) zap.Field { return zap.String(key, input) }
