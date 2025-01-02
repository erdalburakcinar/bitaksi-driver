.PHONY: build run test docker-build docker-run docker-compose-up clean

# Service settings
BINARY_NAME=driver-service
MAIN=cmd/main.go

# Docker settings
DOCKER_IMAGE=driver-service
DOCKER_TAG=latest
DOCKER_COMPOSE_FILE=docker-compose.yml

# Default port
PORT=8080

# Build the application binary
build:
	@echo "Building application..."
	go build -o $(BINARY_NAME) $(MAIN)

# Run the application locally
run: build
	@echo "Running application locally..."
	./$(BINARY_NAME)

# Run tests
test:
	@echo "Running tests..."
	go test ./... -v

# Build the Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

# Run the Docker container
docker-run: docker-build
	@echo "Running Docker container..."
	docker run --rm -p $(PORT):$(PORT) --name $(BINARY_NAME) $(DOCKER_IMAGE):$(DOCKER_TAG)

# Start all services with Docker Compose
docker-compose-up:
	@echo "Starting services with Docker Compose..."
	docker compose -f $(DOCKER_COMPOSE_FILE) up --build

# Stop all services with Docker Compose
docker-compose-down:
	@echo "Stopping services with Docker Compose..."
	docker compose -f $(DOCKER_COMPOSE_FILE) down

# Clean up binaries and Docker images
clean:
	@echo "Cleaning up..."
	rm -f $(BINARY_NAME)
	docker rm -f $(BINARY_NAME) || true
	docker rmi -f $(DOCKER_IMAGE):$(DOCKER_TAG) || true