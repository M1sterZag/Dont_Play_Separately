package core_igdb_provider

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	ClientID     string        `envconfig:"CLIENT_ID" required:"true"`
	ClientSecret string        `envconfig:"CLIENT_SECRET" required:"true"`
	BaseURL      string        `envconfig:"BASE_URL" default:"https://api.igdb.com/v4"`
	Timeout      time.Duration `envconfig:"TIMEOUT" default:"10s"`
}

func NewConfig() (Config, error) {
	var config Config
	if err := envconfig.Process("IGDB", &config); err != nil {
		return Config{}, fmt.Errorf("process config: %w", err)
	}

	return config, nil
}

func NewConfigMust() Config {
	config, err := NewConfig()
	if err != nil {
		err = fmt.Errorf("get igdb config: %w", err)
		panic(err)
	}

	return config
}
