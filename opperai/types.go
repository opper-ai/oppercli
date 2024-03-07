package opperai

// Message represents a single message in a chat.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatPayload is the payload for initiating a chat.
type ChatPayload struct {
	Messages []Message `json:"messages"`
}

// ContextData represents context data.
type ContextData struct {
	IndexID  int                    `json:"index_id"`
	Content  string                 `json:"content"`
	Metadata map[string]interface{} `json:"metadata"` // Use map[string]interface{} to represent a JSON object.
}

// FunctionResponse represents a function response.
type FunctionResponse struct {
	Message     *string       `json:"message"`
	JsonPayload interface{}   `json:"json_payload"` // Use interface{} to represent any JSON value.
	Error       *string       `json:"error"`
	Context     []ContextData `json:"context"`
}

// FunctionDescription represents the description of a function.
type FunctionDescription struct {
	ID           *int                   `json:"id"`
	Path         string                 `json:"path"`
	Description  string                 `json:"description"`
	InputSchema  map[string]interface{} `json:"input_schema"` // Use map[string]interface{} to represent a JSON object.
	OutSchema    map[string]interface{} `json:"out_schema"`   // Use map[string]interface{} to represent a JSON object.
	Instructions string                 `json:"instructions"`
	IndexIDs     []int                  `json:"index_ids"`
}
