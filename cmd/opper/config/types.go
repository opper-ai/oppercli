package config

type APIKeyConfig struct {
	Key string `yaml:"key"`
	// Future fields specific to an API key config can go here
	// BaseURL string `yaml:"base_url"`
}

type Config struct {
	APIKeys map[string]APIKeyConfig `yaml:"api_keys"`
	// Future global config options can go here
	// DefaultModel string `yaml:"default_model"`
}

// Helper method to get the default API key
func (c *Config) GetDefaultAPIKey() string {
	if key, exists := c.APIKeys["default"]; exists {
		return key.Key
	}
	return ""
}
