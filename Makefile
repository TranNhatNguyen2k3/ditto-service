.PHONY: swagger-init swagger-build build run test install-tools lint format init docker-build docker-run docker-stop

MOCK_OUTPUT_DIR=internal/mock
MOCK_CASE=snake

install-tools:
	@echo "Installing required tools..."
	@go install github.com/swaggo/swag/cmd/swag@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install golang.org/x/tools/cmd/goimports@latest
	@go install github.com/vektra/mockery/v2@latest
	@echo "Tools installed successfully!"

swagger-init:
	swag init -g cmd/app/main.go -o docs/swagger

swagger-build:
	swag fmt
	swag init -g cmd/app/main.go -o docs/swagger

init:
	go mod tidy
	go mod download

build:
	go build -o bin/app ./cmd/app/main.go

run:
	bin/app api

test:
	go test -v ./...

fmt:
	goimports -w .

lint:
	golangci-lint run ./...

generate-mocks:
	mockery

docker-build:
	docker-compose build

docker-run:
	docker-compose up -d

docker-stop:
	docker-compose down

clean:
	rm -rf bin/
	rm -rf docs/swagger/
	rm -rf internal/mock/