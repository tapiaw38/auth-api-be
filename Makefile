.DEFAULT_GOAL:=help
.PHONY: help welcome install-deps install-migrate run run-dev clean

welcome: ## Welcome message
	@echo "\033[1;33mIT-AUTH - AUTH API Back-End\033[0m\n"

install-deps: welcome ## Install dependencies for running the project
	@echo "Installing Delve (debugger)..."
	@go install github.com/go-delve/delve/cmd/dlv@latest
	@echo "Delve installation completed"
	@echo "Installing Air (live reloader)..."
	@go install github.com/air-verse/air@latest
	@echo "Air installation completed"

apt-install-migrate: welcome ## Install the migrate CLI tool
	@echo "Installing migrate CLI..."
	@sudo apt install --no-install-recommends -y postgresql-client curl && \
    curl -L -o /tmp/migrate.tar.gz https://github.com/golang-migrate/migrate/releases/download/v4.15.0/migrate.linux-amd64.tar.gz && \
    sudo tar xzf /tmp/migrate.tar.gz -C /usr/local/bin && \
    sudo chmod +x /usr/local/bin/migrate
	@echo "migrate installation completed"

pacman-install-migrate: welcome ## Install the migrate CLI tool
	@echo "Installing migrate CLI..."
	@sudo pacman -S --noconfirm postgresql curl && \
    curl -L -o /tmp/migrate.tar.gz https://github.com/golang-migrate/migrate/releases/download/v4.15.0/migrate.linux-amd64.tar.gz && \
    sudo tar xzf /tmp/migrate.tar.gz -C /usr/local/bin && \
    sudo chmod +x /usr/local/bin/migrate
	@echo "migrate installation completed"

run: welcome ## Run project locally
	@go run ./cmd/api/

run-dev: welcome ## Run project on watch mode for development purposes
	@air

test: welcome ## Run tests project
	@echo "Running tests"
	@go test ./... -cover -coverprofile=coverage.out

cover: welcome test ## Analizes the coverage profiles generated by 'make test' using function mode.
	@go tool cover -func=coverage.out

cover-html: welcome test ## Analizes the coverage profiles generated by 'make test' using html mode.
	@go tool cover -html=coverage.out

gen-mocks: welcome ## Generate all mocks
	@echo "Generating mocks..."
	@go generate ./...

clean: welcome ## Remove unused files and golang cache
	@go clean -cache
	-find . -name "*.out" -exec rm -f  {} \;
	rm -rf bin/

init-docker: welcome ## Initializes a local Docker instance
	@docker-compose up -d

init-mockserver: welcome ## Initializes a local Wiremock instance
	@docker rm mockserver -f
	@docker run -v ./resources/wiremock:/home/wiremock -p 8081:8080 --name mockserver wiremock/wiremock:3.3.1 --verbose

migrate-create: ## Create a new migration
	@read -p "Enter migration name: " name; \
	migrate create -dir migrations -ext sql $$name

migrate-up: welcome ## Apply all available migrations
	@migrate -path migrations -database 'postgres://postgres:postgres@localhost:54321/auth-api-db?sslmode=disable' up

migrate-down: welcome ## Revert the last migration
	@migrate -path migrations -database 'postgres://postgres:postgres@localhost:54321/auth-api-db?sslmode=disable' down

help: welcome
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | grep ^help -v | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'
