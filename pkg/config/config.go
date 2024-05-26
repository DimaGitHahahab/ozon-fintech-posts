package config

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	HTTPPort      string `envconfig:"HTTP_PORT" required:"true"`
	DbURL         string `envconfig:"DB_URL"`
	MigrationPath string `envconfig:"MIGRATION_PATH"`
	Repository    string `envconfig:"REPOSITORY" default:"IN_MEMORY" required:"true"` // IN_MEMORY, POSTGRES
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("failed to load env variables: %w", err)
	}

	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to process env variables: %w", err)
	}

	return &cfg, nil
}
