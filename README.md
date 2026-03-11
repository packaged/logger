# logger

A structured logging package for Go, built on [zap](https://github.com/uber-go/zap). Provides a global singleton logger with environment-aware defaults, common field support, context propagation, timed logging, and GCP Cloud Logging compatibility.

## Install

```bash
go get github.com/packaged/logger/v3
```

## Setup

Initialize the global logger with an environment. The log level and encoding are configured automatically based on the environment.

```go
import (
    "github.com/packaged/environment/environment"
    "github.com/packaged/logger/v3/logger"
)

logger.Setup(environment.Production)
```

| Environment | Encoding | Level |
|---|---|---|
| Production, Sandbox | JSON (GCP format) | Info |
| Integration Test | JSON (GCP format) | Debug |
| Development, Local, Unit Test | Console | Debug |

Debug logging can be forced in any environment by setting the `PACKAGED__DEBUG_LOG=true` environment variable.

## Logging

```go
log := logger.I() // global instance

log.Debug("starting process")
log.Info("user logged in", zap.String("user-id", "abc"))
log.Warn("slow query", zap.Duration("elapsed", elapsed))
log.Error("request failed", zap.Error(err))
```

### Conditional logging

Log only when an error is non-nil:

```go
log.ErrorIf(err, "failed to save", zap.String("id", id))
log.WarnIf(err, "retrying request")
log.InfoIf(err, "completed with warning")
log.DebugIf(err, "optional detail")
```

## Common Fields

Add fields that are included in every log entry from a logger instance. Use `Clone()` to create a scoped logger without mutating the global instance.

```go
log := logger.I().Clone()
log.AddCommon(
    zap.String("request-id", reqID),
    zap.String("service", "payments"),
)

log.Info("processing") // includes request-id and service automatically
```

## Context Propagation

Store and retrieve a logger from `context.Context`, enabling request-scoped loggers to flow through call chains.

```go
// At the request entry point
log := logger.I().Clone()
log.AddCommon(zap.String("request-id", reqID))
ctx = logger.NewContext(ctx, log)

// Anywhere downstream
log := logger.FromContext(ctx)
log.Info("handling step") // includes request-id
```

`FromContext` returns the global logger if none is set on the context, so it is always safe to call.

## Timed Logging

Log at a severity level based on how long an operation took.

```go
tl := logger.DefaultTimedLogConfig().NewLog("db query", zap.String("table", "users"))
defer logger.I().TimedLog(tl)

// ... do work ...
tl.Complete()
```

Default thresholds:

| Duration | Level |
|---|---|
| >= 1 minute | Error |
| >= 30 seconds | Warn |
| >= 2 seconds | Info |
| >= 500ms | Debug |

Custom thresholds:

```go
cfg := &logger.TimedLogConfig{
    ErrorDuration: 10 * time.Second,
    WarnDuration:  5 * time.Second,
    InfoDuration:  time.Second,
    DebugDuration: 100 * time.Millisecond,
}
tl := cfg.NewLog("api call")
```

## Custom Configuration

Create a logger with a custom zap config and options:

```go
l, err := logger.InstanceWithConfig(
    environment.Production,
    zap.NewProductionConfig(),
    logger.WithGoogleEncoding,
    logger.DisableStacktrace,
    logger.Info,
)
```

Available options: `Debug`, `Info`, `Warn`, `Error`, `DPanic`, `Panic`, `Fatal`, `WithConsoleEncoding`, `WithGoogleEncoding`, `DisableStacktrace`, `DisableCaller`.

## Log Data Helpers (`ld` package)

Common zap fields for structured logging:

```go
import "github.com/packaged/logger/v3/ld"

log.Info("request",
    ld.IP(remoteAddr),
    ld.URL(requestURL),
    ld.Method("POST"),
    ld.UserAgent(ua),
    ld.Port(8080),
    ld.Error(err),                        // zap.Skip() if nil
    ld.InterfaceType("handler", h),       // logs the reflect type
    ld.Prefix("req", zap.String("id", id)), // "req:id"
)
```

## Testing

Use `ObserverForTest` to capture log output in tests:

```go
logger.Setup(environment.UnitTest)
logs := logger.ObserverForTest()

logger.I().Info("hello")

assert.Equal(t, 1, logs.Len())
assert.Equal(t, "hello", logs.All()[0].Message)
```
