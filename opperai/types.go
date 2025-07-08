package opperai

import (
	"encoding/json"
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
	Path              string                 `json:"path"`
	Instructions      string                 `json:"instructions"`
	Description       string                 `json:"description"`
	UUID              string                 `json:"uuid"`
	Model             string                 `json:"model"`
	LanguageModelID   int                    `json:"language_model_id"`
	Dataset           Dataset                `json:"dataset"`
	Project           Project                `json:"project"`
	FewShot           bool                   `json:"few_shot"`
	FewShotCount      int                    `json:"few_shot_count"`
	UseSemanticSearch bool                   `json:"use_semantic_search"`
	Revision          int                    `json:"revision"`
	InputSchema       map[string]interface{} `json:"input_schema"`
	OutputSchema      map[string]interface{} `json:"out_schema"`
	IndexConfig       *IndexConfig           `json:"index_config,omitempty"`
	Metadata          map[string]string      `json:"metadata,omitempty"`
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

type Score struct {
	UUID  string  `json:"uuid"`
	Score float64 `json:"score"`
}

type SpanMetric struct {
	SpanUUID string                 `json:"span_uuid"`
	Metrics  map[string]interface{} `json:"metrics"`
}

type Trace struct {
	UUID        string       `json:"uuid"`
	OrgID       int          `json:"org_id"`
	Rating      *int         `json:"rating"`
	Spans       []Span       `json:"spans"`
	Scores      []Score      `json:"scores"`
	StartTime   time.Time    `json:"start_time"`
	EndTime     time.Time    `json:"end_time"`
	DurationMs  float64      `json:"duration_ms"`
	Status      string       `json:"status"`
	Name        string       `json:"name"`
	Input       string       `json:"input"`
	Output      *string      `json:"output"`
	TotalTokens int          `json:"total_tokens"`
	Project     Project      `json:"project"`
	SpanMetrics []SpanMetric `json:"span_metrics"`
}

type Project struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}

type TraceListResponse struct {
	Traces []Trace `json:"traces"`
	Cursor string  `json:"cursor"`
}

type Span struct {
	UUID           string                 `json:"uuid"`
	Project        interface{}            `json:"project"`
	ProjectUUID    *string                `json:"project_uuid"`
	Name           string                 `json:"name"`
	Input          *string                `json:"input"`
	Output         *string                `json:"output"`
	StartTime      time.Time              `json:"start_time"`
	Type           *string                `json:"type"`
	ParentUUID     *string                `json:"parent_uuid"`
	EndTime        time.Time              `json:"end_time"`
	DurationMs     float64                `json:"duration_ms"`
	Error          *string                `json:"error"`
	Meta           map[string]interface{} `json:"meta"`
	Evaluations    interface{}            `json:"evaluations"`
	Score          *float64               `json:"score"`
	TotalTokens    *int                   `json:"total_tokens"`
	DatasetEntries []interface{}          `json:"dataset_entries"`
}

type Dataset struct {
	UUID       string `json:"uuid"`
	EntryCount int    `json:"entry_count"`
}

type BuiltinLanguageModel struct {
	Name            string `json:"name"`
	HostingProvider string `json:"hosting_provider"`
	Location        string `json:"location"`
}

type EvaluationMetric struct {
	Dimension string  `json:"dimension"`
	Value     float64 `json:"value"`
	Comment   string  `json:"comment,omitempty"`
}

type EvaluationStatus struct {
	State   string `json:"state"`
	Details string `json:"details,omitempty"`
}

type EvaluationRecord struct {
	EvaluationRecordUUID string                      `json:"evaluation_record_uuid"`
	DatasetEntryUUID     string                      `json:"dataset_entry_uuid"`
	Status               EvaluationStatus            `json:"status"`
	Input                string                      `json:"input"`
	Expected             string                      `json:"expected"`
	Output               string                      `json:"output"`
	Metrics              map[string]EvaluationMetric `json:"metrics"`
}

type StatisticsSummary struct {
	Sum    float64 `json:"sum"`
	Count  float64 `json:"count"`
	Min    float64 `json:"min"`
	Max    float64 `json:"max"`
	Avg    float64 `json:"avg"`
	Median float64 `json:"median"`
}

type FunctionOverride struct {
	Model        string `json:"model"`
	Instructions string `json:"instructions"`
	FewShotCount *int   `json:"few_shot_count"`
}

type Evaluation struct {
	EvaluationUUID    string                       `json:"evaluation_uuid"`
	DatasetUUID       string                       `json:"dataset_uuid"`
	Records           []EvaluationRecord           `json:"records"`
	Status            EvaluationStatus             `json:"status"`
	Dimensions        []string                     `json:"dimensions"`
	SummaryStatistics map[string]StatisticsSummary `json:"summary_statistics"`
	FunctionOverride  FunctionOverride             `json:"function_override"`
	CreatedAt         string                       `json:"created_at"`
	UpdatedAt         string                       `json:"updated_at"`
}

type EvaluationsResponse struct {
	Meta struct {
		TotalCount int `json:"total_count"`
	} `json:"meta"`
	Data []Evaluation `json:"data"`
}

type UsageStats struct {
	TotalTokensInput  int     `json:"total_tokens_input"`
	TotalTokensOutput int     `json:"total_tokens_output"`
	TotalTokens       int     `json:"total_tokens"`
	TotalCost         float64 `json:"total_cost"`
	Count             int     `json:"count"`
}

type UsageItem struct {
	ProjectName  string    `json:"project_name"`
	FunctionName string    `json:"function_path"`
	Model        string    `json:"model"`
	TokensInput  int       `json:"tokens_input"`
	TokensOutput int       `json:"tokens_output"`
	CreatedAt    time.Time `json:"created_at"`
	Cost         float64   `json:"cost"`
	ID           string    `json:"id"`
	UUID         string    `json:"uuid"`
	TotalTokens  int       `json:"total_tokens"`
}

type UsageEvent struct {
	TimeBucket string                 `json:"time_bucket"`
	Cost       string                 `json:"cost"`
	Count      int                    `json:"count"`
	EventType  *string                `json:"event_type,omitempty"`
	Fields     map[string]interface{} `json:"-"`
}

func (e *UsageEvent) UnmarshalJSON(data []byte) error {
	type Alias UsageEvent
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(e),
	}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	// Create fields map if not exists
	if e.Fields == nil {
		e.Fields = make(map[string]interface{})
	}

	// Unmarshal into a map to get all fields
	var rawMap map[string]interface{}
	if err := json.Unmarshal(data, &rawMap); err != nil {
		return err
	}

	// Add all fields except the standard ones to the Fields map
	for k, v := range rawMap {
		switch k {
		case "time_bucket", "cost", "count", "event_type":
			continue
		default:
			e.Fields[k] = v
		}
	}

	return nil
}

type UsageResponse []UsageEvent

type UsageSummary struct {
	TotalCost            string                 `json:"total_cost"`
	GenerationCost       string                 `json:"generation_cost"`
	PlatformCost         string                 `json:"platform_cost"`
	SpanCost             string                 `json:"span_cost"`
	EmbeddingCost        string                 `json:"embedding_cost"`
	MetricCost           string                 `json:"metric_cost"`
	DatasetStorageCost   string                 `json:"dataset_storage_cost"`
	ImageCost            string                 `json:"image_cost"`
	TotalEvents          int                    `json:"total_events"`
	DateRange            []time.Time            `json:"date_range"`
	EventCountBreakdown  map[string]int         `json:"event_count_breakdown"`
}

type UsageParams struct {
	FromDate    string   `url:"from_date,omitempty"`
	ToDate      string   `url:"to_date,omitempty"`
	EventType   string   `url:"event_type,omitempty"`
	Granularity string   `url:"granularity,omitempty"`
	Fields      []string `url:"fields,omitempty"`
	GroupBy     []string `url:"group_by,comma,omitempty"`
}
