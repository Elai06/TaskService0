BINARY_NAME=taskService

all: build, lint, lint-fix

build: main.go
	go build -o $(BINARY_NAME) main.go

test:
	go test ./...

test-server:
	go test -cover ./api/server

lint:
	golangci-lint run

lint-fix:
	golangci-lint run --fix

clean:
	rm -f $(BINARY_NAME)

run: build
	./$(BINARY_NAME)
