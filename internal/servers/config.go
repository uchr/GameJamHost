package servers

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Config struct {
	HostURI string `env:"HOST_URI,required"`
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
