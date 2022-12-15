module github.com/packaged/logger

go 1.16

require (
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/multierr v1.9.0 // indirect
	go.uber.org/zap v1.24.0
)

// v1.1.0 breaking changes tagged as v1 by mistake
retract v1.1.0
