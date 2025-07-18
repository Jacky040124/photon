package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Jacky040124/photon/pkg"
)

type Config struct {
	OpenRouterKey string `json:"openrouter_key,omitempty"`
	CurrentModel  string `json:"current_model"`
}

// Validate checks if required configuration is present
func (c *Config) Validate() error {
	if c.OpenRouterKey == "" && os.Getenv("PHOTON_OPEN_ROUTER_KEY") == "" {
		return fmt.Errorf("PHOTON_OPEN_ROUTER_KEY environment variable is required")
	}
	// Allow __online__ models to bypass validation
	if c.CurrentModel != "" && !pkg.ValidateModel(c.CurrentModel) && !isOnlineModel(c.CurrentModel) {
		return fmt.Errorf("invalid model '%s'", c.CurrentModel)
	}
	return nil
}

func isOnlineModel(modelID string) bool {
	return len(modelID) > 10 && modelID[:10] == "__online__"
}

// GetOpenRouterKey returns the API key from config or environment
func (c *Config) GetOpenRouterKey() string {
	if c.OpenRouterKey != "" {
		return c.OpenRouterKey
	}
	return os.Getenv("PHOTON_OPEN_ROUTER_KEY")
}

// GetCurrentModel returns the current model, defaulting if not set
func (c *Config) GetCurrentModel() string {
	if c.CurrentModel == "" {
		return pkg.GetDefaultModel()
	}
	return c.CurrentModel
}

// SetCurrentModel updates the current model and saves config
func (c *Config) SetCurrentModel(modelID string) error {
	if !pkg.ValidateModel(modelID) {
		return fmt.Errorf("invalid model '%s'", modelID)
	}

	c.CurrentModel = modelID
	return c.Save()
}

// getConfigPath returns the path to the config file
func getConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configDir := filepath.Join(homeDir, ".photon")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", err
	}

	return filepath.Join(configDir, "config.json"), nil
}

// Save saves the config to disk
func (c *Config) Save() error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

// LoadConfig loads configuration from file and environment
func LoadConfig() (*Config, error) {
	config := &Config{
		CurrentModel: pkg.GetDefaultModel(),
	}

	// Try to load from config file
	configPath, err := getConfigPath()
	if err == nil {
		if data, err := os.ReadFile(configPath); err == nil {
			json.Unmarshal(data, config)
		}
	}

	// Always prefer environment variable for API key
	if envKey := os.Getenv("PHOTON_OPEN_ROUTER_KEY"); envKey != "" {
		config.OpenRouterKey = envKey
	}

	return config, nil
}
