package config

import (
	"context"
	"fmt"
	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	ServerAddress         string `env:"SERVER_ADDRESS,required"`
	ProviderAPIKey        string `env:"PROVIDER_API_KEY,required"`
	ProviderAPIURL        string `env:"PROVIDER_API_URL,required"`
	UpdateTimeInMinutes   int    `env:"UPDATE_TIME_IN_MINUTES,required"`
	ReadTimeoutInSeconds  int    `env:"READ_TIMEOUT_SEC,required"`
	WriteTimeoutInSeconds int    `env:"WRITE_TIMEOUT_SEC, required"`
	DB                    *DBConfig
}

type DBConfig struct {
	DatabaseHost     string `env:"DATABASE_HOST,required"`
	DatabasePort     string `env:"DATABASE_PORT,required"`
	DatabaseName     string `env:"DATABASE_NAME,required"`
	DatabaseUser     string `env:"DATABASE_USER,required"`
	DatabasePassword string `env:"DATABASE_PASS,required"`
}

func New(ctx context.Context) (*Config, error) {
	var cfg Config

	if err := envconfig.Process(ctx, &cfg); err != nil {
		return nil, fmt.Errorf("config process err %w", err)
	}

	return &cfg, nil
}
