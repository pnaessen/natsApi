-include .env

.PHONY:  build run clean lint docker-up docker-down docker-logs deps

APP_NAME := natsApi
MAIN_PATH := ./cmd/api
BINARY_NAME := api
GO := go

BLUE := \033[0;34m
GREEN := \033[0;32m
YELLOW := \033[1;33m
RED := \033[0;31m
NC := \033[0m #


deps:
	${GO} mod download
	${GO} mod verify
	@echo "$(GREEN)✓ OK$(NC)"

build: deps
	@echo "$(BLUE)Compilation...$(NC)"
	${GO} build -o bin/${BINARY_NAME} ${MAIN_PATH}
	@echo "$(GREEN)✓ Build OK - bin/${BINARY_NAME}$(NC)"

run:
	go run ${MAIN_PATH}


lint:
	@command -v golangci-lint >/dev/null 2>&1 || { echo "Init  golangci-lint..."; go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; }
	golangci-lint run

clean:
	@echo "$(BLUE)Clean...$(NC)"
	rm -rf bin/
	${GO} clean
	@echo "$(GREEN)✓ Clean OK$(NC)"

docker-up:
	@echo "$(BLUE)Start de NATS...$(NC)"
	docker run -d --name nats-server -p 4222:4222 -p 8222:8222 nats:latest
	@echo "$(GREEN)✓ NATS localhost:4222$(NC)"

docker-down:
	@echo "$(BLUE)Stop  NATS...$(NC)"
	docker stop nats-server || true
	docker rm nats-server || true
	@echo "$(GREEN)✓ NATS stop$(NC)"

docker-logs:
	@echo "$(BLUE)Logs NATS...$(NC)"
	docker logs -f nats-server

all: clean lint  build

