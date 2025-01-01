package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"bitaksi-go-driver/internal/models"
	"bitaksi-go-driver/internal/repository"
)

type DriverHandler interface {
	ImportLocations(w http.ResponseWriter, r *http.Request)
	FindNearestDriver(w http.ResponseWriter, r *http.Request)
}

type DriverService interface {
	ImportLocations(ctx context.Context) error
	FindNearestDriver(ctx context.Context, latitude, longitude float64, radius int) (*models.DriverWithDistance, error)
}

type driverHandler struct {
	service DriverService
}

func NewDriverHandler(service DriverService) DriverHandler {
	return &driverHandler{service: service}
}

// ImportLocations imports driver locations from a CSV file
// @Summary Import Driver Locations
// @Description Upload driver locations from a predefined CSV file
// @Tags Driver
// @Accept json
// @Produce json
// @Success 200 {string} string "Locations imported successfully"
// @Failure 500 {string} string "Failed to import locations"
// @Router /driver/api/v1/import [post]
func (h *driverHandler) ImportLocations(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	err := h.service.ImportLocations(r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Failed to import locations: %v"}`, err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Locations imported successfully"})
}

// FindNearestDriver searches for nearby drivers within a radius
// @Summary Search Nearby Drivers
// @Description Finds drivers near a given location within the specified radius
// @Tags Driver
// @Param latitude query float64 true "Latitude"
// @Param longitude query float64 true "Longitude"
// @Param radius query int true "Search radius in meters"
// @Success 200 {array} models.DriverWithDistance
// @Failure 400 {string} string "Invalid input"
// @Failure 500 {string} string "Failed to search drivers"
// @Router /driver/api/v1/search [get]
func (h *driverHandler) FindNearestDriver(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Parse query parameters
	latitude, err := strconv.ParseFloat(r.URL.Query().Get("latitude"), 64)
	if err != nil || latitude < -90 || latitude > 90 {
		http.Error(w, `{"error": "Invalid latitude: must be between -90 and 90"}`, http.StatusBadRequest)
		return
	}

	longitude, err := strconv.ParseFloat(r.URL.Query().Get("longitude"), 64)
	if err != nil || longitude < -180 || longitude > 180 {
		http.Error(w, `{"error": "Invalid longitude: must be between -180 and 180"}`, http.StatusBadRequest)
		return
	}

	radius, err := strconv.Atoi(r.URL.Query().Get("radius"))
	if err != nil || radius <= 0 {
		http.Error(w, `{"error": "Invalid radius: must be a positive integer"}`, http.StatusBadRequest)
		return
	}

	// Call the service
	results, err := h.service.FindNearestDriver(r.Context(), latitude, longitude, radius)
	if err != nil {
		if errors.Is(err, repository.ErrDriverNotFound) {
			http.Error(w, `{"error": "No drivers found"}`, http.StatusNotFound)
			return
		}

		http.Error(w, fmt.Sprintf(`{"error": "Failed to search drivers: %v"}`, err), http.StatusInternalServerError)
		return
	}

	// Send response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(results)
}
