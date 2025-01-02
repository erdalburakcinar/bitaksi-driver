# Driver Service

Driver Service is a microservice responsible for managing driver locations and providing APIs to search for the nearest drivers. It interacts with MongoDB for data storage and supports Docker and Docker Compose for containerization and orchestration.

---

## Features

- **Swagger Documentation**: Interactive API documentation.
- **Dockerized**: Fully containerized for deployment.
- **Authentication**: Secured endpoints with API keys.

---

## Requirements

- [Go](https://golang.org/doc/install) 1.23
- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)
- MongoDB (containerized via Docker Compose)

---

## Setup

### Build the Application
```bash
make build
```
Run the Application Locally

```bash
make run
```

Run Tests

```bash
make test
```

Docker Setup

Use Docker Compose
```bash
make docker-compose-up
```

Stop All Services

```bash	
make docker-compose-down
```
API Documentation

Driver Service includes Swagger for API documentation. After running the service, visit the following URL in your browser:

http://localhost:8080/swagger/index.html

This provides an interactive interface to test and explore all available endpoints.

Example Usage with Swagger
	1.	Open your browser and navigate to http://localhost:8080/swagger/index.html.
	2.	Authenticate using your API key if required.
	3.	Explore the endpoints, test API calls, and view detailed response structures.

