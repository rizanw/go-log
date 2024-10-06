package log

import (
	"context"

	"github.com/google/uuid"
)

const (
	KeyCtxRequestID = "request_id"
	KeyCtxUserInfo  = "user_info"
	KeyCtxSource    = "source"
)

// SetCtxRequestID generates & sets request_id to context
func SetCtxRequestID(ctx context.Context, requestID string) context.Context {
	if requestID != "" {
		return context.WithValue(ctx, KeyCtxRequestID, requestID)
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

func SetCtxUserInfo(ctx context.Context, userInfo interface{}) context.Context {
	if userInfo != nil {
		return context.WithValue(ctx, KeyCtxUserInfo, userInfo)
	}
	return ctx
}

func GetCtxUserInfo(ctx context.Context) interface{} {
	if ctx == nil {
		return nil
	}

	if userInfo, ok := ctx.Value(KeyCtxUserInfo).(interface{}); ok {
		return userInfo
	}
	return nil
}

func SetCtxSource(ctx context.Context, source interface{}) context.Context {
	if source != nil {
		return context.WithValue(ctx, KeyCtxSource, source)
	}
	return ctx
}

func GetCtxSource(ctx context.Context) interface{} {
	if ctx == nil {
		return nil
	}

	if source, ok := ctx.Value(KeyCtxSource).(interface{}); ok {
		return source
	}
	return nil
}
