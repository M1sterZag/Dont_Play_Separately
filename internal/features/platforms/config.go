package platforms_config

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	CacheTTL time.Duration `envconfig:"CACHE_TTL" default:"1h"`
}

func NewConfig() (Config, error) {
	var config Config
	if err := envconfig.Process("PLATFORMS", &config); err != nil {
		return Config{}, fmt.Errorf("process config: %w", err)
	}

	return config, nil
}

func NewConfigMust() Config {
	config, err := NewConfig()
	if err != nil {
		err = fmt.Errorf("get platforms config: %w", err)
		panic(err)
	}

	return config
}