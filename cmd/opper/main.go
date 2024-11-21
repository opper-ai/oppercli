package main

import (
	"bufio"
	"context"
	"flag"
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

func getAPIKey(keyName string) (string, error) {
	// First check environment variable for specific key name
	if envKeyName := os.Getenv("OPPER_KEY_NAME"); envKeyName != "" {
		keyName = envKeyName
	}

	// Check environment variable for API key
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
		if apiKey := cfg.GetAPIKey(keyName); apiKey != "" {
			return apiKey, nil
		}
	}

	// If no key found, prompt user
	fmt.Printf("No API key found in environment variable OPPER_API_KEY or config file ~/.oppercli")
	if keyName != "" {
		fmt.Printf(" for key '%s'", keyName)
	}
	fmt.Println()

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

	fmt.Printf("Enter your API key for '%s': ", keyName)
	apiKey, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("error reading API key: %w", err)
	}
	apiKey = strings.TrimSpace(apiKey)

	// Create new config with API key
	if cfg == nil {
		cfg = &config.Config{
			APIKeys: make(map[string]config.APIKeyConfig),
		}
	}
	cfg.APIKeys[keyName] = config.APIKeyConfig{
		Key: apiKey,
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
	// Parse global flags first
	var keyName string
	flag.StringVar(&keyName, "key", "", "Name of the API key to use")
	flag.Parse()

	// Get remaining args
	args := flag.Args()
	if len(args) == 0 {
		args = []string{"help"}
	}

	// Add version flag handling
	if len(args) > 0 && (args[0] == "--version" || args[0] == "-v") {
		fmt.Printf("opper version %s\n", Version)
		os.Exit(0)
	}

	// Get API key from environment or config file
	apiKey, err := getAPIKey(keyName)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	// Get base URL from environment or config file
	baseURL := os.Getenv("OPPER_BASE_URL")
	if baseURL == "" {
		if cfg, err := getConfig(); err == nil && cfg != nil {
			baseURL = cfg.GetBaseUrl(keyName)
		}
	}

	// Initialize the client
	client := opperai.NewClient(apiKey, baseURL)

	// Create command parser
	parser := commands.NewCommandParser()

	// Parse command using the remaining args
	// Add program name back to args for parser
	fullArgs := append([]string{"opper"}, args...)
	cmd, err := parser.Parse(fullArgs)
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

// Add helper function to get config
func getConfig() (*config.Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("could not get home directory: %w", err)
	}

	configPath := filepath.Join(homeDir, ".oppercli")
	return readConfig(configPath)
}
