package auth_config

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	JWTSecret     string        `envconfig:"SECRET" required:"true"`
	JWTAccessTTL  time.Duration `envconfig:"ACCESS_TTL" default:"1h"`
	JWTRefreshTTL time.Duration `envconfig:"REFRESH_TTL" default:"168h"`
}

func NewConfig() (Config, error) {
	var config Config
	if err := envconfig.Process("JWT", &config); err != nil {
		return Config{}, fmt.Errorf("process config: %w", err)
	}

	return config, nil
}

func NewConfigMust() Config {
	config, err := NewConfig()
	if err != nil {
		err = fmt.Errorf("get jwt config: %w", err)
		panic(err)
	}

	return config
}
