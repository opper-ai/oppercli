package config

type APIKeyConfig struct {
	Key string `yaml:"key"`
}

type Config struct {
	APIKeys map[string]APIKeyConfig `yaml:"api_keys"`
}

// GetAPIKey returns the API key for the given name, or default if name is empty
func (c *Config) GetAPIKey(name string) string {
	if name == "" {
		name = "default"
	}
	if key, exists := c.APIKeys[name]; exists {
		return key.Key
	}
	return ""
}
