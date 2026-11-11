package log

import (
	"context"

	"github.com/google/uuid"
)

// CtxKey is the type of the context keys used by this package. A dedicated
// type prevents collisions with plain string keys set by other packages.
type CtxKey string

const (
	KeyCtxRequestID CtxKey = "request_id"
	KeyCtxUserInfo  CtxKey = "user_info"
	KeyCtxSource    CtxKey = "source"
)

// SetCtxRequestID generates & sets request_id to context
func SetCtxRequestID(ctx context.Context, requestID ...string) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	if len(requestID) > 0 && requestID[0] != "" {
		return context.WithValue(ctx, KeyCtxRequestID, requestID[0])
	}

	return context.WithValue(ctx, KeyCtxRequestID, "gen-"+uuid.New().String())
}

// GetCtxRequestID returns request_id from context
func GetCtxRequestID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	if requestID, ok := ctx.Value(KeyCtxRequestID).(string); ok {
		return requestID
	}
	return ""
}

// SetCtxUserInfo stores user info (e.g. log.KV{"user_id": 42}) to be logged as `user_info`
func SetCtxUserInfo(ctx context.Context, userInfo interface{}) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	if userInfo != nil {
		return context.WithValue(ctx, KeyCtxUserInfo, userInfo)
	}
	return ctx
}

// GetCtxUserInfo returns user info from context
func GetCtxUserInfo(ctx context.Context) interface{} {
	if ctx == nil {
		return nil
	}

	return ctx.Value(KeyCtxUserInfo)
}

// SetCtxSource stores the request source (e.g. client app and version) to be logged as `source`
func SetCtxSource(ctx context.Context, source interface{}) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	if source != nil {
		return context.WithValue(ctx, KeyCtxSource, source)
	}
	return ctx
}

// GetCtxSource returns source from context
func GetCtxSource(ctx context.Context) interface{} {
	if ctx == nil {
		return nil
	}

	return ctx.Value(KeyCtxSource)
}
