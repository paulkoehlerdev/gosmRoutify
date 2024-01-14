.PHONY: build

build: build-importer build-router

build-importer:
	@echo "Building importer..."
	@go build -o bin/importer ./cmd/importer

build-router:
	@echo "Building router..."
	@go build -o bin/router ./cmd/router