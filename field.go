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
		fields.UserInfo = GetCtxUserInfo(ctx)
	}

	if len(metadata) > 0 {
		fields.Metadata = metadata
	}

	return fields
}
