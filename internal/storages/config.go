package storages

import (
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Config struct {
	ConnectTimeout time.Duration `env:"CONNECT_TIMEOUT,required"`
	DatabasePath   string        `env:"DATABASE_PATH,required"`

	MigrationPath    string `env:"MIGRATION_PATH"`
	MigrationVersion int32  `env:"MIGRATION_VERSION"`
	MigrationEnabled bool   `env:"MIGRATION_ENABLED"`
}

func NewConfig() (*Config, error) {
	var err = godotenv.Load()
	if err != nil {
		return nil, err
	}

	cfg := Config{}
	err = env.Parse(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
