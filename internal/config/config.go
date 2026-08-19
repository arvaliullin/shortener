package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	ServerAddress string `envconfig:"SERVER_ADDRESS" default:":8080"`
	BaseURL       string `envconfig:"BASE_URL" default:"http://localhost:8080"`
	RedisAddr     string `envconfig:"REDIS_ADDR" required:"true"`
}

func Load() (*Config, error) {
	var c Config
	if err := envconfig.Process("", &c); err != nil {
		return nil, fmt.Errorf("ошибка чтения переменных окружения: %w", err)
	}
	return &c, nil
}
