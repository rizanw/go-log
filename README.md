# go-log

[![Go Reference](https://pkg.go.dev/badge/github.com/rizanw/go-log.svg)](https://pkg.go.dev/github.com/rizanw/go-log)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

A structured, leveled logging package for Go services. It wraps [zerolog](https://github.com/rs/zerolog) or [zap](https://github.com/uber-go/zap) behind one small API, adds request-scoped fields from `context.Context` (request ID, user, source), and masks sensitive values such as passwords and tokens before they reach your logs.

```go
log.Info(ctx, nil, log.KV{"order_id": 1042}, "order created")
```

```json
{"level":"info","app":"checkout","env":"production","request_id":"5825511e-196f-406b-baed-67a9da40a26a","metadata":{"order_id":1042},"timestamp":"2024-10-23T14:52:00Z","message":"order created"}
```

## Features

- **Two engines, one API**: switch between zerolog (default) and zap with a config value
- **Leveled logging**: `Debug`, `Info`, `Warn`, `Error`, `Fatal`, each with a `printf`-style variant
- **Context fields**: request ID, user info, and source are read from `ctx` and added to every log line
- **Sensitive data masking**: values under keys like `password` or `token` are masked at any depth, without modifying your data
- **Output options**: console or JSON, optional colors, log file, or file and console together
- **Caller and stack traces**: optional file:line and stack traces for errors

## Installation

```bash
go get github.com/rizanw/go-log
```

Requires Go 1.20 or newer.

## Quick start

It works without any setup:

```go
import (
    "context"

    log "github.com/rizanw/go-log"
)

func main() {
    log.Infof(context.Background(), nil, log.KV{"port": 8080}, "starting %s", "go-app")
}
```

For services, configure it once at startup:

```go
err := log.SetConfig(&log.Config{
    AppName:           "go-app",
    Environment:       "production",
    Level:             log.InfoLevel,
    UseJSON:           true,
    WithCaller:        true,
    MaskSensitiveData: []string{"password", "token", "authorization"},
})
if err != nil {
    log.Fatal(context.Background(), err, nil, "invalid log config")
}
```

Call `SetConfig` before logging from other goroutines. The zerolog engine applies some settings (such as time format) globally.

## Logging

Every function takes the same leading arguments:

```go
log.Error(ctx, err, log.KV{"user_id": 42}, "failed to charge card")
//        │    │    │                       └ message
//        │    │    └ metadata: extra fields (nil if none)
//        │    └ error to attach (nil if none)
//        └ context carrying request-scoped fields
```

| Level | Plain | Formatted | Use for |
|---|---|---|---|
| Debug | `log.Debug` | `log.Debugf` | detailed diagnostics while developing or troubleshooting |
| Info | `log.Info` | `log.Infof` | normal milestones: startup, completed jobs, config changes |
| Warn | `log.Warn` | `log.Warnf` | something unexpected that the app recovered from |
| Error | `log.Error` | `log.Errorf` | a failed operation that needs attention |
| Fatal | `log.Fatal` | `log.Fatalf` | an unrecoverable error; **logs and then exits the process** |

Lines below the configured `Level` are dropped.

## Context fields

Set request-scoped values once, for example in HTTP middleware, and they appear on every log line that uses that context:

```go
ctx = log.SetCtxRequestID(ctx)            // generates a request ID
ctx = log.SetCtxRequestID(ctx, requestID) // or use your own, e.g. from an X-Request-ID header
ctx = log.SetCtxUserInfo(ctx, log.KV{"user_id": 42, "role": "admin"})
ctx = log.SetCtxSource(ctx, log.KV{"app": "ios", "version": "1.10.5"})

log.Info(ctx, nil, nil, "profile updated")
```

```json
{"level":"info","request_id":"gen-3f1c…","user_info":{"role":"admin","user_id":42},"source":{"app":"ios","version":"1.10.5"},"message":"profile updated"}
```

`GetCtxRequestID`, `GetCtxUserInfo`, and `GetCtxSource` read the values back, for example to forward the request ID to another service.

## Masking sensitive data

List the keys to mask in `MaskSensitiveData`. Matching is case-insensitive and works at any depth in metadata and user info, including nested maps and slices, and for values of any type:

```go
log.SetConfig(&log.Config{UseJSON: true, MaskSensitiveData: []string{"password", "pin"}})

req := log.KV{"username": "ann", "password": "hunter2", "card": log.KV{"PIN": 1234}}
log.Info(ctx, nil, log.KV{"request": req}, "login attempt")
// {"metadata":{"request":{"card":{"PIN":"****"},"password":"*******","username":"ann"}}, ...}
```

The masked values are written to a copy, so `req` itself still holds the real password after logging. To change how values are masked, set `SensitiveDataMasker`:

```go
SensitiveDataMasker: func(value string) string { return "[REDACTED]" },
```

## Configuration

| Field | Type | Default | Description |
|---|---|---|---|
| `AppName` | `string` | | logged as `app` |
| `Environment` | `string` | | logged as `env`; `dev`, `development`, or `local` enables development mode |
| `Level` | `log.Level` | `DebugLevel` | minimum level to log |
| `TimeFormat` | `string` | RFC 3339 | timestamp format |
| `UseJSON` | `bool` | `false` | JSON output instead of human-readable console output |
| `UseColor` | `bool` | `false` | colored console output (zerolog engine, development mode only) |
| `WithCaller` | `bool` | `false` | add the `file:line` that logged the message |
| `CallerSkip` | `int` | `0` | extra stack frames to skip, if you wrap this package |
| `WithStack` | `bool` | `false` | add stack traces for errors |
| `StackLevel` | `*log.Level` | `ErrorLevel` | minimum level for stack traces (zap engine) |
| `StackMarshaller` | `func(error) interface{}` | `pkgerrors.MarshalStack` | how stack traces are extracted (zerolog engine) |
| `MaskSensitiveData` | `[]string` | | keys whose values are masked |
| `SensitiveDataMasker` | `func(string) string` | `*` per character | how masked values are rendered |
| `FilePath` | `string` | | write logs to this file instead of the console (directories are created) |
| `UseMultiWriters` | `bool` | `false` | write to both `FilePath` and stdout (requires `FilePath`) |
| `Engine` | `log.Engine` | `log.Zerolog` | `log.Zerolog` or `log.Zap` |

Capturing the caller or a stack trace costs extra allocations on every log call, so enable them where the information is worth it.

### Choosing an engine

Both are fast, structured, zero-allocation loggers, and the output is equivalent. zerolog is the default and has the nicer colored console output for local development. Pick zap if your codebase already uses it. For a comparison, see [Better Stack's guide to Go logging libraries](https://betterstack.com/community/guides/logging/best-golang-logging-libraries/).

## Example

[`example/main.go`](example/main.go) shows configuration, error logging with stack traces, context fields, masking, and file plus console output:

```bash
go run ./example
```

## License

[MIT](LICENSE)
