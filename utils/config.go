package utils

import (
	"time"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	// Database
	DatabaseURL string `env:"DATABASE_URL,required"`

	// Server
	Port     string `env:"PORT" envDefault:"8080"`
	LogLevel string `env:"LOG_LEVEL" envDefault:"info"`

	// API Security
	APIKey string `env:"API_KEY,required"`

	// External Services
	WOM struct {
		BaseURL string        `env:"WOM_BASE_URL" envDefault:"https://api.wiseoldman.net/v2"`
		Timeout time.Duration `env:"WOM_TIMEOUT" envDefault:"30s"`
	}

	// Railway detection (for logging format)
	RailwayProjectID string `env:"RAILWAY_PROJECT_ID"`

	// Optional features
	AllowedOrigins []string `env:"ALLOWED_ORIGINS" envSeparator:"," envDefault:"*"`
}

func LoadConfig() (*Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
