SERVICE=bmt-saas
MAIN=./cmd/server
BIN=./bin/$(SERVICE)

.PHONY: build dev test test-race test-integration lint sqlc migrate-up migrate-down infra-up infra-down tidy

build:
	go build -o $(BIN) $(MAIN)

dev:
	go run $(MAIN)/main.go

test:
	go test ./... -v -timeout 60s

test-race:
	go test -race ./... -timeout 60s

test-integration:
	go test ./... -v -timeout 120s -tags integration

lint:
	golangci-lint run ./...

sqlc:
	sqlc generate

migrate-up:
	migrate -path migrations -database "$$(grep DATABASE_URL .env | cut -d= -f2-)" up

migrate-down:
	migrate -path migrations -database "$$(grep DATABASE_URL .env | cut -d= -f2-)" down 1

migrate-create:
	migrate create -ext sql -dir migrations -seq $(name)

infra-up:
	docker compose up -d

infra-down:
	docker compose down

tidy:
	go mod tidy

.PHONY: run
run: build
	$(BIN)
