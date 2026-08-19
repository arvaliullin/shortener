.PHONY: run
run:
	set -a && \
	. ./.env && \
	set +a && \
	go run github.com/arvaliullin/shortener/cmd/shortener

.PHONY: fmt
fmt:
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
generate:
	- go generate ./...

.PHONY: up
up:
	docker compose -f deployments/docker-compose.yaml up -d --build

.PHONY: down
down:
	docker compose -f deployments/docker-compose.yaml down -v

.PHONY: ps
ps:
	docker compose -f deployments/docker-compose.yaml ps --orphans=false
.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: test
test:
	- go test ./...
