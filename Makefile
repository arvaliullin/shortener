.DEFAULT_GOAL := help

.PHONY: help
help: ## Показать список команд
	@awk 'BEGIN {FS = ":.*## "} /^[a-zA-Z0-9_-]+:.*## / {printf "  %-16s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: run
run: ## Запустить API-сервер
	set -a && \
	. ./.env && \
	set +a && \
	go run github.com/arvaliullin/shortener/cmd/shortener

.PHONY: web-install
web-install: ## Установить зависимости фронтенда
	npm --prefix web/shortener ci

.PHONY: web-dev
web-dev: certs web-install ## Запустить фронтенд по HTTPS
	npm --prefix web/shortener run dev

.PHONY: sync
sync: ## Установить зависимости Python
	uv sync

.PHONY: fmt
fmt: ## Отформатировать Go-код
	- go fmt ./...

.PHONY: install-deps
install-deps: ## Установить инструменты разработки
	- go install go.uber.org/mock/mockgen@v0.6.0
	- go install github.com/pressly/goose/v3/cmd/goose@latest
	- go install github.com/swaggo/swag/cmd/swag@latest
	- go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
	- go get github.com/go-chi/chi/v5
	- go get github.com/redis/go-redis/v9
	- go get github.com/kelseyhightower/envconfig
	- go get github.com/rs/zerolog

.PHONY: generate
generate: ## Сгенерировать код
	- go generate ./...

.PHONY: certs
certs: ## Выписать локальный TLS-сертификат
	@command -v mkcert >/dev/null 2>&1 || { printf '%s\n' "mkcert не найден. Установите mkcert и повторите make certs."; exit 1; }
	mkdir -p certs
	mkcert -cert-file certs/localhost.pem -key-file certs/localhost-key.pem localhost 127.0.0.1 ::1

.PHONY: up
up: certs ## Собрать и запустить стек в Docker
	docker compose -f deployments/docker-compose.yaml up -d --build

.PHONY: down
down: ## Остановить стек и удалить тома
	docker compose -f deployments/docker-compose.yaml down -v

.PHONY: ps
ps: ## Показать контейнеры стека
	docker compose -f deployments/docker-compose.yaml ps --orphans=false

.PHONY: logs
logs: ## Показать журналы web и shortener
	docker compose -f deployments/docker-compose.yaml logs -f web shortener

.PHONY: lint
lint: ## Проверить Go-код линтером
	golangci-lint run ./...

.PHONY: test
test: ## Запустить тесты
	- go test ./...

.PHONY: check
check: generate fmt lint test ## Сгенерировать код, отформатировать, проверить и протестировать
