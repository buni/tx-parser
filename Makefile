up:
	docker compose up --build -d 
down: 
	docker compose down
cleanup:
	docker compose down -v --rmi all --remove-orphans
generate: 
	go generate ./...
lint: 
	go tool golangci-lint run ./...
