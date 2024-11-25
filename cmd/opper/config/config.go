package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// LoadConfig reads the configuration from ~/.oppercli
func LoadConfig() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("error getting home directory: %w", err)
	}

	configPath := filepath.Join(home, ".oppercli")
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{APIKeys: make(map[string]APIKeyConfig)}, nil
		}
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("error parsing config file: %w", err)
	}

	if config.APIKeys == nil {
		config.APIKeys = make(map[string]APIKeyConfig)
	}

	return &config, nil
}

// SaveConfig writes the configuration to ~/.oppercli
func SaveConfig(config *Config) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("error getting home directory: %w", err)
	}

	configPath := filepath.Join(home, ".oppercli")
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("error marshaling config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return fmt.Errorf("error writing config file: %w", err)
	}

	return nil
}

// GetAPIKey returns the API key from environment or config file
func GetAPIKey(key string) (string, error) {
	// First check environment variable
	envKey := os.Getenv("OPPER_API_KEY")
	if envKey != "" {
		return envKey, nil
	}

	// Then check config file
	config, err := LoadConfig()
	if err != nil {
		return "", err
	}

	if apiKeyConfig, ok := config.APIKeys[key]; ok {
		return apiKeyConfig.Key, nil
	}

	// Prompt user to save API key
	fmt.Printf("No API key found in environment variable OPPER_API_KEY or config file ~/.oppercli for key '%s'\n", key)
	fmt.Print("Would you like to save an API key now? (y/n): ")
	var response string
	fmt.Scanln(&response)
	if response == "y" || response == "Y" {
		fmt.Print("Enter API key: ")
		var apiKey string
		fmt.Scanln(&apiKey)

		if config.APIKeys == nil {
			config.APIKeys = make(map[string]APIKeyConfig)
		}
		config.APIKeys[key] = APIKeyConfig{Key: apiKey}
		if err := SaveConfig(config); err != nil {
			return "", fmt.Errorf("error saving config: %w", err)
		}
		return apiKey, nil
	}

	return "", fmt.Errorf("no API key found")
}

// GetAPIKeyAndBaseUrl returns both the API key and base URL
func GetAPIKeyAndBaseUrl(key string) (string, string, error) {
	// First check environment variable
	envKey := os.Getenv("OPPER_API_KEY")
	if envKey != "" {
		baseUrl := os.Getenv("OPPER_BASE_URL")
		return envKey, baseUrl, nil
	}

	// Then check config file
	config, err := LoadConfig()
	if err != nil {
		return "", "", err
	}

	if apiKeyConfig, ok := config.APIKeys[key]; ok {
		return apiKeyConfig.Key, apiKeyConfig.BaseUrl, nil
	}

	return "", "", fmt.Errorf("no API key found")
}

// ValidateConfig checks if the configuration is valid
func ValidateConfig(config *Config) error {
	if config.APIKeys == nil {
		return fmt.Errorf("api_keys section is required")
	}

	for name, key := range config.APIKeys {
		if key.Key == "" {
			return fmt.Errorf("API key for '%s' is empty", name)
		}
		if key.BaseUrl != "" {
			// Optional: Add URL validation here
			if !strings.HasPrefix(key.BaseUrl, "http://") && !strings.HasPrefix(key.BaseUrl, "https://") {
				return fmt.Errorf("invalid base URL for '%s': must start with http:// or https://", name)
			}
		}
	}
	return nil
}
