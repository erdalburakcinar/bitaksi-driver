# Use Go 1.23 as the base image
FROM golang:1.23-alpine as builder

WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./
RUN go mod download

# Copy the application source code
COPY . .

# Build the application binary
RUN go build -o driver-service ./cmd/main.go

# Use a minimal base image for running the application
FROM alpine:latest

WORKDIR /app

# Copy the built binary and required files
COPY --from=builder /app/driver-service .
COPY --from=builder /app/docs/Coordinates.csv ./docs/Coordinates.csv
COPY --from=builder /app/config/config.yaml ./config/config.yaml

# Expose the application port
EXPOSE 8080

CMD ["./driver-service"]