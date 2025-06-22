APP_NAME=simple-product-api

.PHONY: run build test test-product test-cover test-bench-product wire mocks-all mocks-product lint mocks coverage docker-up docker-down fmt swag

run:
	go run ./cmd/main.go

build:
	go build -o bin/$(APP_NAME) ./cmd/main.go

test:
	go test ./... -

test-product:
	go test ./internal/product/usecase -v

test-cover:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out

test-bench-product:
	go test -bench=. -benchmem ./internal/product/usecase

wire :
	wire gen ./pkg/di

lint:
	golangci-lint run

fmt:
	go fmt ./...

mocks-all:
	mockery --all --output=internal/mocks

mocks-product:
	mockery --name=ProductRepository --output=usecase/mocks --with-expecter

coverage-md:
	go-cover-markdown < coverage.out > COVERAGE.md

swag:
	swag init -g cmd/main.go -o docs

docker-up:
	docker-compose up --build

docker-down:
	docker-compose down
