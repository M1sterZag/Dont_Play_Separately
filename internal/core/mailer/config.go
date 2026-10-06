package core_mailer

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Host     string `envconfig:"HOST" required:"true"`
	Port     int    `envconfig:"PORT" default:"25"`
	User     string `envconfig:"USER"`
	Password string `envconfig:"PASSWORD"`
	From     string `envconfig:"FROM" required:"true"`
}

func NewConfig() (Config, error) {
	var config Config
	if err := envconfig.Process("SMTP", &config); err != nil {
		return Config{}, fmt.Errorf("process smtp config: %w", err)
	}

	return config, nil
}

func NewConfigMust() Config {
	config, err := NewConfig()
	if err != nil {
		err = fmt.Errorf("get smtp config: %w", err)
		panic(err)
	}

	return config
}