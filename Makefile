BINARY_NAME=taskService

all: build, lint

build: main.go
	go build -o $(BINARY_NAME) main.go

test:
	go test ./...

lint:
	golangci-lint run --fix

clean:
	rm -f $(BINARY_NAME)

run: build
	./$(BINARY_NAME)
