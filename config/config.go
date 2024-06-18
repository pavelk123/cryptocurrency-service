package config

import (
	"context"
	"fmt"

	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	ServerAddress         string         `env:"SERVER_ADDRESS,required"`
	UpdateTimeInMinutes   int            `env:"UPDATE_TIME_IN_MINUTES,required"`
	ReadTimeoutInSeconds  int            `env:"READ_TIMEOUT_SEC,default=10"`
	WriteTimeoutInSeconds int            `env:"WRITE_TIMEOUT_SEC, default=10"`
	DB                    DBConfig       `env:",prefix=DATABASE_"`
	Provider              ProviderConfig `env:",prefix=PROVIDER_"`
}

type DBConfig struct {
	Host     string `env:"HOST,required"`
	Port     string `env:"PORT,required"`
	Name     string `env:"NAME,required"`
	User     string `env:"USER,required"`
	Password string `env:"PASS,required"`
}

type ProviderConfig struct {
	URL       string   `env:"API_URL,required"`
	Key       string   `env:"API_KEY,required"`
	Currecies []string `env:"CURRENCIES"`
}

func New(ctx context.Context) (*Config, error) {
	var cfg Config

	if err := envconfig.Process(ctx, &cfg); err != nil {
		return nil, fmt.Errorf("config process err: %w", err)
	}

	return &cfg, nil
}
