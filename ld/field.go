// Package ld Log Data - quick zap fields
package ld

// Log Data

import (
	"reflect"

	"go.uber.org/zap"
)

// Prefix prepends a prefix to the field key, separated by a colon.  This is useful for grouping fields together.
func Prefix(prefix string, field zap.Field) zap.Field {
	field.Key = prefix + ":" + field.Key
	return field
}

// Error returns a zap.Field for an error, or zap.Skip() if the error is nil.
func Error(err error) zap.Field {
	if err == nil {
		return zap.Skip()
	}
	return zap.Error(err)
}

// InterfaceType returns a zap.Field for the type of the interface.
func InterfaceType(key string, iface any) zap.Field {
	return zap.String(key+":type", reflect.TypeOf(iface).String())
}

// IP returns a zap.Field for an IP address.
func IP(ip string) zap.Field { return zap.String("ip", ip) }

// UserAgent returns a zap.Field for a user agent.
func UserAgent(ua string) zap.Field { return zap.String("user-agent", ua) }

// URL returns a zap.Field for a URL.
func URL(input string) zap.Field { return zap.String("url", input) }

// PortString returns a zap.Field for a port.
func PortString(input string) zap.Field { return zap.String("port", input) }

// Port returns a zap.Field for a port.
func Port(input int) zap.Field { return zap.Int("port", input) }

// Method returns a zap.Field for a method.
func Method(input string) zap.Field { return zap.String("method", input) }
