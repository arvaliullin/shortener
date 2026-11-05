// Package config читает настройки процесса агента.
package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

// Config задает адрес агента, модель и соседние сервисы.
type Config struct {
	Address          string `envconfig:"AGENT_ADDRESS" default:":8090"`
	OllamaBaseURL    string `envconfig:"OLLAMA_BASE_URL" default:"http://127.0.0.1:11434/v1"`
	OllamaModel      string `envconfig:"OLLAMA_MODEL" default:"qwen2.5:3b"`
	MCPURL           string `envconfig:"MCP_URL" default:"http://127.0.0.1:8080/mcp"`
	KeycloakIssuer   string `envconfig:"KEYCLOAK_ISSUER" default:"https://localhost:8443/realms/shortener"`
	KeycloakJWKSURL  string `envconfig:"KEYCLOAK_JWKS_URL" default:"http://127.0.0.1:8082/realms/shortener/protocol/openid-connect/certs"`
	KeycloakClientID string `envconfig:"KEYCLOAK_CLIENT_ID" default:"shortener-web"`
}

// Load читает настройки из переменных окружения.
func Load() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("read environment: %w", err)
	}
	return &cfg, nil
}
