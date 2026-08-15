package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	HTTPAddr  string `envconfig:"HTTP_ADDR" default:":8080"`
	RedisAddr string `envconfig:"REDIS_ADDR" required:"true"`
}

func Load() (*Config, error) {
	var c Config
	if err := envconfig.Process("", &c); err != nil {
		return nil, fmt.Errorf("ошибка чтения переменных окружения: %w", err)
	}
	return &c, nil
}
