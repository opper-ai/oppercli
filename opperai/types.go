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
	Path         string            `json:"path"`
	Instructions string            `json:"instructions"`
	Description  string            `json:"description"`
	Model        string            `json:"model,omitempty"`        // For custom model support
	IndexConfig  *IndexConfig      `json:"index_config,omitempty"` // For index support
	Metadata     map[string]string `json:"metadata,omitempty"`     // For additional configuration
}

type IndexConfig struct {
	Type       string            `json:"type"`       // e.g., "vector", "keyword"
	Source     string            `json:"source"`     // Data source for the index
	Settings   map[string]string `json:"settings"`   // Index-specific settings
	Dimensions int               `json:"dimensions"` // For vector indexes
}

type CustomLanguageModel struct {
	Name           string                 `json:"name"`
	Identifier     string                 `json:"identifier"`
	APIKey         string                 `json:"api_key"`
	Extra          map[string]interface{} `json:"extra"`
	ID             int                    `json:"id"`
	OrganizationID int                    `json:"organization_id"`
	CreatedAt      string                 `json:"created_at"`
	UpdatedAt      string                 `json:"updated_at"`
}

// Function represents an AI function
type Function struct {
	Path         string `json:"path"`
	Description  string `json:"description,omitempty"`
	Instructions string `json:"instructions"`
}
