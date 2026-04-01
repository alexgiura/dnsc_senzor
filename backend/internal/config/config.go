package config

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type AppSettings struct {
	ServerPort  string `env:"SERVER_PORT" envDefault:"8080"`
	Environment string `env:"ENVIRONMENT" envDefault:"development"`
	DebugMode   bool   `env:"DEBUG_MODE" envDefault:"false"`
}

type NetworkAlertsConfig struct {
	StoragePath string `env:"NETWORK_ALERTS_STORAGE_PATH" envDefault:"data/network_alerts.jsonl"`
}

type Config struct {
	AppSettings   AppSettings
	NetworkAlerts NetworkAlertsConfig
}

func Load() (*Config, error) {
	cfg := &Config{}

	// Get the current file path
	_, currentFilePath, _, _ := runtime.Caller(0)

	// Navigate to backend directory:
	backendPath := filepath.Join(filepath.Dir(currentFilePath), "..", "..", "..")
	envFilePath := filepath.Join(backendPath, ".env")

	// Load the .env file from backend directory
	err := godotenv.Load(envFilePath)
	if err != nil {
		log.Printf("No .env file found at %s, using environment variables.\n", envFilePath)
	} else {
		log.Printf("✅ Loaded .env file from: %s\n", envFilePath)
	}

	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("error loading configuration: %s", err)
	}

	if cfg.AppSettings.ServerPort == "" {
		return nil, fmt.Errorf("invalid config: SERVER_PORT must not be empty")
	}

	log.Printf("✅ Loaded config - ServerPort: %s, Environment: %s\n",
		cfg.AppSettings.ServerPort, cfg.AppSettings.Environment)

	return cfg, nil
}
