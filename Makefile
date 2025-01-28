DOCKER_COMPOSE_TOOLS= docker compose up -d && docker compose exec tools

up:
	docker compose up --build -d 
down: 
	docker compose down
cleanup:
	docker compose down -v --rmi all --remove-orphans
generate: ## generate mocks and types from openapi spec
	go install go.uber.org/mock/mockgen@v0.5.0 
	go install github.com/dmarkham/enumer@v1.5.9
	go generate ./...
lint: ## generate mocks and types from openapi spec
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.62.2 
	golangci-lint run ./...