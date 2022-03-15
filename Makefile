GOPATH:=$(shell go env GOPATH)

.PHONY: build
build:
	go build -o build/main cmd/main.go

.PHONY: format
format:
	gofmt -s -w .

.PHONY: lint
lint:
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.44.2
	golangci-lint run ./...
