package service

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strconv"

	"bitaksi-go-driver/internal/models"
	"bitaksi-go-driver/internal/repository"
)

const (
	filePath = "./docs/Coordinates.csv" // Define the constant file path
)

// DriverRepository provides methods to interact with driver data
type DriverRepository interface {
	SaveDrivers(ctx context.Context, locations []models.DriverWithDistance) error
	FindNearestDriver(ctx context.Context, latitude, longitude float64, radius int) (*models.DriverWithDistance, error)
	EnsureIndex(ctx context.Context) error
}

type DriverService struct {
	repo DriverRepository
}

func NewDriverService(repo DriverRepository) DriverService {
	return DriverService{repo: repo}
}

func (s *DriverService) ImportLocations(ctx context.Context) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open CSV file at %s: %w", filePath, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read CSV file: %w", err)
	}

	var locations []models.DriverWithDistance
	for i, record := range records {
		if i == 0 {
			continue // Skip header
		}

		latitude, err1 := strconv.ParseFloat(record[0], 64)
		longitude, err2 := strconv.ParseFloat(record[1], 64)
		if err1 != nil || err2 != nil {
			return fmt.Errorf("invalid data format in CSV at row %d", i+1)
		}

		locations = append(locations, models.DriverWithDistance{
			Location: models.Location{
				Type:        "Point",
				Coordinates: []float64{longitude, latitude},
			},
		})
	}

	if err := s.repo.SaveDrivers(ctx, locations); err != nil {
		return fmt.Errorf("failed to save drivers to repository: %w", err)
	}

	return nil
}

func (s *DriverService) FindNearestDriver(ctx context.Context, latitude, longitude float64, radius int) (*models.DriverWithDistance, error) {
	// Ensure the geospatial index exists
	if err := s.repo.EnsureIndex(ctx); err != nil {
		return nil, fmt.Errorf("failed to ensure index: %w", err)
	}

	// Perform the geospatial search
	driver, err := s.repo.FindNearestDriver(ctx, latitude, longitude, radius)
	if err != nil {
		if errors.Is(err, repository.ErrDriverNotFound) {
			return nil, fmt.Errorf("no drivers found within the radius of %d meters: %w", radius, err)
		}
		return nil, fmt.Errorf("failed to find nearest driver: %w", err)
	}

	return driver, nil
}
