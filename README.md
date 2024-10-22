# go-log

go-log is a structured logging package to provides a simple way to log messages in different levels.
It is designed to be simple and easy to use.

## Quick Start

```bash
go get -u github.com/rizanw/go-log
```

And you can drop this package in to replace your log very simply:

```go
package main

import (
	"context"

	// import this log package
	log "github.com/rizanw/go-log"
)

func main() {
	var (
		appName = "go-app"
	)

	// use the log like this
	log.Infof(context.TODO(), nil, log.KV{"key": "value"}, "starting app %s", appName)
}
```

## Configurable Log

You can also customize configuration for your log very simple:

```go
package main

import (
	"context"
	// import this log package
	log "github.com/rizanw/go-log"
)

func main() {
	var (
		ctx = context.Background()
		err error
	)

	// customize configuration like this
	err = log.SetConfig(&log.Config{
		AppName:     "go-app",
		Environment: "development",
		WithCaller:  true,
	})
	if err != nil {
		// use the log like this
		log.Errorf(ctx, err, nil, "we got error with message: %s", err.Error())
	}
}

```

Note: `SetConfig` func is not thread safe, so call it once when initializing the app only.

### Configuration

below is list of available configuration:

| Key             | type                        | Description                                                                        |
|-----------------|-----------------------------|------------------------------------------------------------------------------------|
| AppName         | string                      | your application name                                                              |
| Environment     | string                      | your application environment                                                       |
| Level           | log.Level                   | minimum log level to be printed (default: DEBUG)                                   |
| TimeFormat      | string                      | desired time format (default: RFC3339)                                             |
| WithCaller      | bool                        | caller toggle to print which line is calling the log (default: false)              |
| CallerSkip      | int                         | which caller line wants to be print                                                |
| WithStack       | bool                        | toggle to print which stack trace error located (default: false)                   |
| StackLevel      | log.Level                   | minimum log level for zap stack trace (default: ERROR)                             |
| StackMarshaller | func(err error) interface{} | function to get and log the stack trace for zerolog (default: `zerolog/pkgerrors`) |
| UseMultiWriters | bool                        | a toggle to print log into log file and log console (FilePath required)            |
| FilePath        | string                      | specify your output log files directories (default: no file)                       |
| UseJSON         | bool                        | a toggle to format log as json (default: false)                                    |
| UseColor        | bool                        | a toggle to colorize your log console with zerolog                                 |
| Engine          | log.Engine                  | desired engine logger (default: zerolog)                                           |                      

note:

- Keep in mind that taking a caller or stacktrace is eager and expensive (relatively speaking) and makes an additional
  allocation.

#### Engine Options

This pkg currently provides two engine (aka logger) to use:

- [Zerolog](https://github.com/rs/zerolog)
- [Zap](https://github.com/uber-go/zap)

if you confused to decide, you can
read [this article](https://betterstack.com/community/guides/logging/best-golang-logging-libraries/) as reference.

## Structured Log

by implementing structured logging, we can easily filter and search logs based on the key-value fields:

```json
{
  "level": "info",
  "timestamp": "2024-07-23T14:52:00Z",
  "app": "golang-app",
  "env": "development",
  "request_id": "5825511e-196f-406b-baed-67a9da40a26a",
  "source": {
    "app": "ios",
    "version": "1.10.5"
  },
  "metadata": {
    "username": "hello",
    "password": "***"
  },
  "message": "[HTTP][Request]: POST /api/v1/login"
}
```

## Hierarchical Log

this package provide 5 hierarchical levels based on the severity:

- **DEBUG** - this log level is used to obtain diagnostic information that can be helpful for troubleshooting and
  debugging. These messages often contain verbose or fine-grained information about the inner workings of the system or
  application. When teams look for log data to filter out for cost savings, they often start with DEBUG logs.
- **INFO** - this log level provide general information about the status of the system. This log level is useful for
  tracking an
  application's progress or operational milestones. For example, your application may create INFO logs upon application
  startup, when a user makes configuration changes, or when they successfully complete tasks.
- **WARN** - this log level serves as a signal for potential issues that are not necessarily a critical error. For
  example, your system may generate a WARN log when it is short on resources. If WARN logs go unaddressed, they may lead
  to bigger issues in the future.
- **ERROR** - this log level indicates significant problems that happened in the system. It usually denotes that an
  unexpected event or exception has occurred. This log level is crucial for identifying issues affecting user experience
  or overall functionality, so immediate attention is needed.
- **FATAL** - this log level shows severe conditions that cause the system to terminate or operate in a significantly
  degraded state. These logs are used for serious problems, like crashes or conditions that threaten data integrity or
  application stability. FATAL logs often lead to service disruptions.

## Usage

### logging

you can log based on the severity hierarchy and each severity has 2 main functions:

- unformatted, similar to `Println` in `fmt` package

```go
// Debug
log.Debug(ctx, err, log.KV{}, "this is a debug log")
// Info
log.Info(ctx, err, log.KV{}, "this is an info log")
// Warn
log.Warn(ctx, err, log.KV{}, "this is a warning log")
// Error
log.Error(ctx, err, log.KV{}, "this is an error log")
// Fatal
log.Fatal(ctx, err, log.KV{}, "this is a fatal log")
```

- formatted, similar to `Printf` in `fmt` package

```go
// Debug
log.Debugf(ctx, err, log.KV{}, "this is a debug log: %s", err.Error())
// Info
log.Infof(ctx, err, log.KV{}, "this is an info log: %s", err.Error())
// Warn
log.Warnf(ctx, err, log.KV{}, "this is a warning log: %s", err.Error())
// Error
log.Errorf(ctx, err, log.KV{}, "this is an error log: %s", err.Error())
// Fatal
log.Fatalf(ctx, err, log.KV{}, "this is a fatal log: %s", err.Error())
```

### context

```go
// set request_id (generated) logging into context
ctx = log.SetRequestID(ctx)

// set request_id (by yours) logging into context
ctx = log.SetRequestID(ctx, requestID)
```

```go
// set user_info logging into context
ctx = log.SetUserInfo(ctx, log.KV{"username": "hello"})
```

```go
// set source logging into context
ctx = log.SetSource(ctx, log.KV{"app": source.App, "version": source.Version})
```

### Additional Fields

Need more fields? coming soon!
