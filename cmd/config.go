package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type Config struct {
	OpenRouterKey string
}

// Validate checks if required configuration is present
func (c *Config) Validate() error {
	if c.OpenRouterKey == "" {
		return fmt.Errorf("OPEN_ROUTER_KEY environment variable is required")
	}
	return nil
}

func LoadConfig() (*Config, error) {
	envPath := filepath.Join("..", "configs", ".env")
	if _, err := os.Stat(envPath); err == nil {
		godotenv.Load(envPath)
	}

	config := &Config{
		OpenRouterKey: os.Getenv("OPEN_ROUTER_KEY"),
	}

	return config, nil
}
