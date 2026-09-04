package core_storage

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Endpoint      string `envconfig:"ENDPOINT" default:"localhost:9000"`
	AccessKey     string `envconfig:"ACCESS_KEY" required:"true"`
	SecretKey     string `envconfig:"SECRET_KEY" required:"true"`
	Bucket        string `envconfig:"BUCKET" default:"avatars"`
	PublicBaseURL string `envconfig:"PUBLIC_BASE_URL" default:"http://localhost:9000"`
	UseSSL        bool   `envconfig:"USE_SSL" default:"false"`
	Region        string `envconfig:"REGION" default:"us-east-1"`
}

func NewConfig() (Config, error) {
	var config Config
	if err := envconfig.Process("S3", &config); err != nil {
		return Config{}, fmt.Errorf("process config: %w", err)
	}

	return config, nil
}

func NewConfigMust() Config {
	config, err := NewConfig()
	if err != nil {
		err = fmt.Errorf("get s3 config: %w", err)
		panic(err)
	}

	return config
}
