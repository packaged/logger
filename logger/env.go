// Package logger is a helper package for logging
package logger

import "github.com/packaged/environment/environment"

const (
	// BinaryDebugLogging is the environment variable that can be used to enable debug logging in a binary
	BinaryDebugLogging environment.Name = "PACKAGED__DEBUG_LOG"
)
