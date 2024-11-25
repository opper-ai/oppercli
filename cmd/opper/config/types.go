package config

type APIKeyConfig struct {
	Key     string `yaml:"key"`
	BaseUrl string `yaml:"baseUrl,omitempty"`
}

type Config struct {
	APIKeys map[string]APIKeyConfig `yaml:"api_keys"`
}
