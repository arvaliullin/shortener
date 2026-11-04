package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	ServerAddress    string `envconfig:"SERVER_ADDRESS" default:":8080"`
	BaseURL          string `envconfig:"BASE_URL" default:"http://localhost:8080"`
	RedisAddr        string `envconfig:"REDIS_ADDR" required:"true"`
	KeycloakIssuer   string `envconfig:"KEYCLOAK_ISSUER" default:"https://localhost:8443/realms/shortener"`
	KeycloakJWKSURL  string `envconfig:"KEYCLOAK_JWKS_URL" default:"http://127.0.0.1:8082/realms/shortener/protocol/openid-connect/certs"`
	KeycloakClientID string `envconfig:"KEYCLOAK_CLIENT_ID" default:"shortener-web"`
}

func Load() (*Config, error) {
	var c Config
	if err := envconfig.Process("", &c); err != nil {
		return nil, fmt.Errorf("ошибка чтения переменных окружения: %w", err)
	}
	return &c, nil
}
