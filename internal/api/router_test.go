package api

import (
	"bitaksi-go-driver/internal/config"
	"bitaksi-go-driver/internal/models"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

// MockDriverService is a mock implementation of the DriverService interface
type MockDriverService struct{}

func (m *MockDriverService) ImportLocations(ctx context.Context) error {
	return nil
}

func (m *MockDriverService) FindNearestDriver(ctx context.Context, latitude, longitude float64, radius int) (*models.DriverWithDistance, error) {
	return nil, nil
}

func TestSetupRouter(t *testing.T) {
	cfg := &config.Config{
		Server: struct {
			APIKey string `mapstructure:"api_key"`
		}{
			APIKey: "test-api-key",
		},
	}

	driverService := &MockDriverService{}
	router := SetupRouter(driverService, cfg)

	tests := []struct {
		name           string
		method         string
		endpoint       string
		headers        map[string]string
		expectedStatus int
	}{
		{
			name:           "Health Check",
			method:         http.MethodGet,
			endpoint:       "/health",
			headers:        nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Unauthorized Import",
			method:         http.MethodPost,
			endpoint:       "/driver/api/v1/import",
			headers:        nil,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Authorized Import",
			method:         http.MethodPost,
			endpoint:       "/driver/api/v1/import",
			headers:        map[string]string{"Authorization": "test-api-key"},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Unauthorized Search",
			method:         http.MethodGet,
			endpoint:       "/driver/api/v1/search?latitude=40.0&longitude=-73.0&radius=5000",
			headers:        nil,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Authorized Search",
			method:         http.MethodGet,
			endpoint:       "/driver/api/v1/search?latitude=40.0&longitude=-73.0&radius=5000",
			headers:        map[string]string{"Authorization": "test-api-key"},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.endpoint, nil)
			if tt.headers != nil {
				for key, value := range tt.headers {
					req.Header.Set(key, value)
				}
			}

			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rec.Code)
			}
		})
	}
}
