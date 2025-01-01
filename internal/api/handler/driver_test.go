package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"bitaksi-go-driver/internal/models"
	"bitaksi-go-driver/internal/repository"
)

type MockDriverService struct {
	ImportLocationsFn   func(ctx context.Context) error
	FindNearestDriverFn func(ctx context.Context, latitude, longitude float64, radius int) (*models.DriverWithDistance, error)
}

func (m *MockDriverService) ImportLocations(ctx context.Context) error {
	return m.ImportLocationsFn(ctx)
}

func (m *MockDriverService) FindNearestDriver(ctx context.Context, latitude, longitude float64, radius int) (*models.DriverWithDistance, error) {
	return m.FindNearestDriverFn(ctx, latitude, longitude, radius)
}

func TestImportLocations(t *testing.T) {
	mockService := &MockDriverService{
		ImportLocationsFn: func(ctx context.Context) error {
			return nil
		},
	}

	handler := NewDriverHandler(mockService)

	req := httptest.NewRequest(http.MethodPost, "/import-locations", nil)
	rec := httptest.NewRecorder()

	handler.ImportLocations(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var response map[string]string
	err := json.NewDecoder(rec.Body).Decode(&response)
	if err != nil || response["message"] != "Locations imported successfully" {
		t.Fatalf("unexpected response: %v", response)
	}
}

func TestFindNearestDriver_TableDriven(t *testing.T) {
	mockService := &MockDriverService{}
	predefinedID, _ := primitive.ObjectIDFromHex("6775be842e9ffeeae6b1de93")

	tests := []struct {
		name           string
		latitude       string
		longitude      string
		radius         string
		mockResponse   *models.DriverWithDistance
		mockError      error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Valid Request",
			latitude:       "40.748817",
			longitude:      "-73.985428",
			radius:         "5000",
			mockResponse:   &models.DriverWithDistance{ID: predefinedID, Distance: 500},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id":"6775be842e9ffeeae6b1de93","location":{"type":"","coordinates":null},"distance":500}`},
		{
			name:           "No Drivers Found",
			latitude:       "40.748817",
			longitude:      "-73.985428",
			radius:         "5000",
			mockResponse:   nil,
			mockError:      repository.ErrDriverNotFound,
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error": "No drivers found"}`,
		},
		{
			name:           "Invalid Latitude",
			latitude:       "invalid",
			longitude:      "-73.985428",
			radius:         "5000",
			mockResponse:   nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error": "Invalid latitude: must be between -90 and 90"}`,
		},
		{
			name:           "Invalid Longitude",
			latitude:       "40.748817",
			longitude:      "invalid",
			radius:         "5000",
			mockResponse:   nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error": "Invalid longitude: must be between -180 and 180"}`,
		},
		{
			name:           "Invalid Radius",
			latitude:       "40.748817",
			longitude:      "-73.985428",
			radius:         "-1",
			mockResponse:   nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error": "Invalid radius: must be a positive integer"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService = &MockDriverService{
				FindNearestDriverFn: func(ctx context.Context, latitude, longitude float64, radius int) (*models.DriverWithDistance, error) {
					return tt.mockResponse, tt.mockError
				},
			}

			handler := NewDriverHandler(mockService)

			req := httptest.NewRequest(http.MethodGet, "/search-nearby", nil)
			query := req.URL.Query()
			query.Add("latitude", tt.latitude)
			query.Add("longitude", tt.longitude)
			query.Add("radius", tt.radius)
			req.URL.RawQuery = query.Encode()

			rec := httptest.NewRecorder()

			handler.FindNearestDriver(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rec.Code)
			}

			if rec.Body.String() != tt.expectedBody+"\n" {
				t.Errorf("expected body %q, got %q", tt.expectedBody, rec.Body.String())
			}
		})
	}
}
