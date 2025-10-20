PORT ?= 8083

run:
	go run ./cmd/server

tidy:
	go mod tidy

test:
	go test ./... -v

build:
	go build -o bin/server ./cmd/server

clean:
	rm -rf bin

help:
	@echo "make run [PORT=8081] - Run server"
	@echo "make build           - Build binary"
	@echo "make tidy            - Go mod tidy"
	@echo "make test            - Run tests"
	@echo "make clean           - Remove bin/"