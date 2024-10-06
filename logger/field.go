package logger

const (
	FieldNameRequestID = "request_id"
	FieldNameMetadata  = "metadata"
)

type Field struct {
	RequestID string
	Metadata  map[string]interface{}
	Fields    map[string]interface{}
}
