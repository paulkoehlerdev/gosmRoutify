.PHONY: build

build: build-loader build-router

build-loader:
	@echo "Building loader..."
	@CGO_ENABLED=1 go build -tags fts5,json -o bin/loader ./cmd/loader/main.go

build-router:
	@echo "Building router..."
	@CGO_ENABLED=1 go build -tags fts5,json -o bin/router ./cmd/router/main.go