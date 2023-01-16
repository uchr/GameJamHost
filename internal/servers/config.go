package servers

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Config struct {
	HostURI    string `env:"HOST_URI,required"`
	ImageDir   string `env:"IMAGE_DIR,required"`
	StaticDir  string `env:"STATIC_DIR,required"`
	SessionKey string `env:"SESSION_KEY,required"`
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
