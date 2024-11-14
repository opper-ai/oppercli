package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/opper-ai/oppercli/cmd/opper/commands"
	"github.com/opper-ai/oppercli/cmd/opper/config"
	"github.com/opper-ai/oppercli/opperai"
)

var Version = "dev"

func getAPIKey() (string, error) {
	// First check environment variable
	if apiKey := os.Getenv("OPPER_API_KEY"); apiKey != "" {
		return apiKey, nil
	}

	// Check config file
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not get home directory: %w", err)
	}

	configPath := filepath.Join(homeDir, ".oppercli")
	cfg, err := readConfig(configPath)
	if err == nil {
		if apiKey := cfg.GetDefaultAPIKey(); apiKey != "" {
			return apiKey, nil
		}
	}

	// If no key found, prompt user
	fmt.Println("No API key found in environment variable OPPER_API_KEY or config file ~/.oppercli")
	fmt.Print("Would you like to save an API key now? (y/n): ")
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("error reading input: %w", err)
	}

	response = strings.TrimSpace(strings.ToLower(response))
	if response != "y" && response != "yes" {
		return "", fmt.Errorf("API key is required to use opper CLI")
	}

	fmt.Print("Enter your API key: ")
	apiKey, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("error reading API key: %w", err)
	}
	apiKey = strings.TrimSpace(apiKey)

	// Create new config with API key
	cfg = &config.Config{
		APIKeys: map[string]config.APIKeyConfig{
			"default": {
				Key: apiKey,
			},
		},
	}

	// Save to config file
	err = saveConfig(configPath, cfg)
	if err != nil {
		return "", fmt.Errorf("error saving config: %w", err)
	}

	fmt.Printf("API key saved to %s\n", configPath)
	return apiKey, nil
}

func readConfig(path string) (*config.Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg config.Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, fmt.Errorf("error parsing config file: %w", err)
	}

	// Initialize the map if it's nil
	if cfg.APIKeys == nil {
		cfg.APIKeys = make(map[string]config.APIKeyConfig)
	}

	return &cfg, nil
}

func saveConfig(path string, cfg *config.Config) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("error marshaling config: %w", err)
	}

	return os.WriteFile(path, data, 0600)
}

func main() {
	// Add version flag handling
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		fmt.Printf("opper version %s\n", Version)
		os.Exit(0)
	}

	// Get API key from environment or config file
	apiKey, err := getAPIKey()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	// Initialize the client with base URL from environment or empty string
	baseURL := os.Getenv("OPPER_BASE_URL")
	client := opperai.NewClient(apiKey, baseURL)

	// Create command parser
	parser := commands.NewCommandParser()

	// Parse command
	cmd, err := parser.Parse(os.Args)
	if err != nil {
		fmt.Println("Error parsing command:", err)
		os.Exit(1)
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Execute command
	if err := cmd.Execute(ctx, client); err != nil {
		fmt.Println("Error executing command:", err)
		os.Exit(1)
	}
}
