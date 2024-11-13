package opperai

import (
	"time"
)

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

type Index struct {
	ID        int       `json:"id"`
	UUID      string    `json:"uuid"`
	Name      string    `json:"name"`
	Files     []File    `json:"files"`
	CreatedAt time.Time `json:"created_at"`
}

type File struct {
	ID               int       `json:"id"`
	OriginalFilename string    `json:"original_filename"`
	Size             int64     `json:"size"`
	IndexStatus      string    `json:"index_status"`
	Key              string    `json:"key"`
	UUID             string    `json:"uuid"`
	CreatedAt        time.Time `json:"created_at"`
}

type Document struct {
	Key      string                 `json:"key"`
	Content  string                 `json:"content"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

type RetrievalResponse struct {
	Key      string                 `json:"key"`
	Content  string                 `json:"content"`
	Score    float64                `json:"score"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

type Filter struct {
	Field     string      `json:"field"`
	Operation string      `json:"operation"`
	Value     interface{} `json:"value"`
}
