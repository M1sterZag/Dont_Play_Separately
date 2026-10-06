package auth_config

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	JWTSecret               string        `envconfig:"SECRET" required:"true"`
	JWTAccessTTL            time.Duration `envconfig:"ACCESS_TTL" default:"1h"`
	JWTRefreshTTL           time.Duration `envconfig:"REFRESH_TTL" default:"168h"`
	VerificationCodeTTL     time.Duration `envconfig:"VERIFICATION_CODE_TTL" default:"10m"`
	VerificationCodeLength  int           `envconfig:"VERIFICATION_CODE_LENGTH" default:"6"`
	MaxVerificationAttempts int           `envconfig:"MAX_VERIFICATION_ATTEMPTS" default:"5"`
	MinResendInterval       time.Duration `envconfig:"MIN_RESEND_INTERVAL" default:"60s"`
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
