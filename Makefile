.PHONY: build

build: build-loader build-router

build-loader:
	@echo "Building loader..."
	@go build -o bin/loader ./cmd/loader

build-router:
	@echo "Building router..."
	@go build -o bin/router ./cmd/router