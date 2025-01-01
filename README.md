# Driver Service

## Overview

The Driver Service is a backend service built using Go, MongoDB, and the Gorilla Mux framework. It provides APIs to manage driver locations, search for nearby drivers, and handle geospatial queries with MongoDB.

---

## Features

- **Import Driver Locations**: Bulk upload of driver locations from a CSV file.
- **Search Nearby Drivers**: Find the nearest drivers within a specified radius from a given location.
- **Geospatial Indexing**: Automatically ensures geospatial indexing for MongoDB collections.

---

## Prerequisites

- **Go**: Version 1.20 or higher
- **Docker**: Installed and running
- **MongoDB**: Accessible instance for database operations

---

## Getting Started

Run with Docker

Build and Run the Service
make docker-build
make docker-run

The service will be accessible at:
http://localhost:8080