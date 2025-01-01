.PHONY: build run test docker-build docker-run clean

# Go application
BINARY_NAME=driver-service
MAIN=cmd/main.go

# Docker settings
DOCKER_IMAGE=driver-service
DOCKER_TAG=latest

# Default port
PORT=8080

build:
	@echo "Building application..."
	go build -o $(BINARY_NAME) $(MAIN)

run: build
	@echo "Running application..."
	./$(BINARY_NAME)

test:
	@echo "Running tests..."
	go test ./... -v

docker-build:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

docker-run:
	@echo "Running Docker container..."
	docker run --rm -d -p $(PORT):$(PORT) --name $(BINARY_NAME) $(DOCKER_IMAGE):$(DOCKER_TAG)

clean:
	@echo "Cleaning up..."
	rm -f $(BINARY_NAME)
	docker rm -f $(BINARY_NAME) || true
	docker rmi -f $(DOCKER_IMAGE):$(DOCKER_TAG) || true