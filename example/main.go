package main

import (
	"context"

	"github.com/pkg/errors"
	"github.com/rizanw/go-log"
)

func main() {
	ctx := context.Background()

	// example log without config
	logWithoutConfig(ctx)

	// example set log config
	setConfig()

	// example log error
	err := errors.New("error-test")
	logError(ctx, err)

	// example log with data
	req := map[string]interface{}{
		"input":    []string{"value-1", "value-2"},
		"password": "secret",
	}
	res := map[string]interface{}{
		"output": "value",
	}
	logWithData(ctx, req, res)

	// example log error with stopping entire app
	logFatal(ctx, err)
}

// logWithoutConfig is example when using go-log without SetConfig
func logWithoutConfig(ctx context.Context) {
	// 2024-10-22T18:58:03+07:00 INF sample log without config
	log.Info(ctx, nil, nil, "sample log without config")
}

// setConfig is example when you need to configure your log
// note: do once when starting your app
func setConfig() {
	err := log.SetConfig(&log.Config{
		AppName:           "go-app",
		Environment:       "local",
		WithCaller:        true,
		WithStack:         true,
		UseJSON:           true,
		UseMultiWriters:   true,
		MaskSensitiveData: []string{"password"},
		FilePath:          "/Users/rizanw/go/src/github.com/rizanw/go-log/example/file.log",
	})
	if err != nil {
		log.Fatal(context.TODO(), err, nil, "failed: set log config")
		return
	}
}

// logError is example to log an error
// note: to add stacktrace, don't forget to set
func logError(ctx context.Context, err error) {
	// {"level":"error","app":"go-app","env":"local","stacktrace":[{"func":"main","line":"20","source":"main.go"},{"func":"main","line":"271","source":"proc.go"},{"func":"goexit","line":"1222","source":"asm_arm64.s"}],"error":"error-test","timestamp":"2024-10-22T19:05:16+07:00","caller":"~/go/src/github.com/rizanw/go-log/example/main.go:61","message":"sample error"}
	log.Error(ctx, err, nil, "sample error")
}

// logWithData is example for log with metadata fields
func logWithData(ctx context.Context, req interface{}, res interface{}) {
	// example adding request_id
	ctx = log.SetCtxRequestID(ctx)

	// {"level":"info","app":"go-app","env":"local","metadata":{"request":{"input":["value-1","value-2"]},"response":{"output":"value"}},"request_id":"gen-ceff830d-4bcb-4dbd-9c51-50e38f57f34d","timestamp":"2024-10-22T19:10:19+07:00","caller
	//":"~/go/src/github.com/rizanw/go-log/example/main.go:69","message":"sample log with data"}
	log.Info(ctx, nil, log.KV{"request": req, "response": res}, "sample log with data")
}

// logFatal is example for log an error with stopping the entire process
func logFatal(ctx context.Context, err error) {
	// {"level":"fatal","app":"go-app","env":"local","error":"error-test","timestamp":"2024-10-22T18:58:03+07:00","caller":"~/go/src/github.com/rizanw/go-log/example/main.go:69","message":"log fatal will also kill your app"}
	log.Fatal(ctx, err, nil, "log fatal will also kill your app")
}
