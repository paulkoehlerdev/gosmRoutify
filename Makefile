
setup:
	cp config.json.example config.json

run:
	go run cmd/gosmRoutify/main.go -config config.json
