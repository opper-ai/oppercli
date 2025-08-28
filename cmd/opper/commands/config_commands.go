package commands

import (
	"context"
	"fmt"

	"github.com/opper-ai/oppercli/cmd/opper/config"
	"github.com/opper-ai/oppercli/opperai"
)

func (c *ConfigCommand) Execute(ctx context.Context, client *opperai.Client) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	switch c.Action {
	case "list":
		fmt.Println("Configured API keys:")
		for name, key := range cfg.APIKeys {
			if key.BaseUrl != "" {
				fmt.Printf("  %s: %s (baseUrl: %s)\n", name, truncateString(key.Key, 10), key.BaseUrl)
			} else {
				fmt.Printf("  %s: %s\n", name, truncateString(key.Key, 10))
			}
		}

	case "add":
		if c.Name == "" || c.Key == "" {
			return fmt.Errorf("name and API key required")
		}
		cfg.APIKeys[c.Name] = config.APIKeyConfig{
			Key:     c.Key,
			BaseUrl: c.BaseUrl,
		}
		if err := config.SaveConfig(cfg); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}
		fmt.Printf("Added API key '%s'\n", c.Name)

	case "remove":
		if _, exists := cfg.APIKeys[c.Name]; !exists {
			return fmt.Errorf("API key '%s' not found", c.Name)
		}
		delete(cfg.APIKeys, c.Name)
		if err := config.SaveConfig(cfg); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}
		fmt.Printf("Removed API key '%s'\n", c.Name)

	case "get":
		if c.Name == "" {
			return fmt.Errorf("name required")
		}
		if apiKey, exists := cfg.APIKeys[c.Name]; exists {
			fmt.Print(apiKey.Key)
		} else {
			return fmt.Errorf("API key '%s' not found", c.Name)
		}

	default:
		return fmt.Errorf("unknown config action: %s", c.Action)
	}

	return nil
}
