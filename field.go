package log

import (
	"context"

	"github.com/rizanw/go-log/logger"
)

type KV map[string]interface{}

func buildFields(ctx context.Context, metadata KV) logger.Field {
	var fields logger.Field

	if ctx != nil {
		fields.RequestID = GetCtxRequestID(ctx)
		fields.Source = GetCtxSource(ctx)

		userInfo := GetCtxUserInfo(ctx)
		if ui, ok := userInfo.(KV); ok {
			// assert alias into origin type
			var m map[string]interface{} = ui
			userInfo = m
		}
		fields.UserInfo = userInfo
	}

	if len(metadata) > 0 {
		fields.Metadata = metadata
	}

	return fields
}
