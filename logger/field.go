package logger

const (
	FieldNameRequestID = "request_id"
	FieldNameSource    = "source"
	FieldNameUserInfo  = "user_info"
	FieldNameMetadata  = "metadata"
)

type Field struct {
	RequestID string
	Source    interface{}
	UserInfo  interface{}
	Metadata  map[string]interface{}
	Fields    map[string]interface{}
}
