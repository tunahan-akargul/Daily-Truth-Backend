run:
	go run ./cmd/main.go

tidy:
	go mod tidy

test:
	go test ./... -v

build:
	go build -o bin/main ./cmd/main.go

clean:
	rm -rf bin

help:
	@echo "make run    - Run the app"
	@echo "make tidy   - Sync deps"
	@echo "make test   - Run tests"
	@echo "make build  - Build binary"
	@echo "make clean  - Remove build artifacts"