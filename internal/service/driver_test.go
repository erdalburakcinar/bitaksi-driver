package service

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/mock"

	"bitaksi-go-driver/internal/models"
	"bitaksi-go-driver/internal/repository"
)

// MockDriverRepository is a mocked implementation of the DriverRepository
type MockDriverRepository struct {
	mock.Mock
}

func (m *MockDriverRepository) SaveDrivers(ctx context.Context, locations []models.DriverWithDistance) error {
	args := m.Called(ctx, locations)
	return args.Error(0)
}

func (m *MockDriverRepository) EnsureIndex(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockDriverRepository) FindNearestDriver(ctx context.Context, latitude, longitude float64, radius int) (*models.DriverWithDistance, error) {
	args := m.Called(ctx, latitude, longitude, radius)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.DriverWithDistance), args.Error(1)
}

func createTestCSVFile(t *testing.T, content string) string {
	t.Helper()

	// Create the test CSV file
	filePath := "./docs/Coordinates.csv"
	err := os.MkdirAll("./docs", os.ModePerm)
	if err != nil {
		t.Fatalf("failed to create docs directory: %v", err)
	}

	err = os.WriteFile(filePath, []byte(content), os.ModePerm)
	if err != nil {
		t.Fatalf("failed to create test CSV file: %v", err)
	}

	return filePath
}

func removeTestCSVFile(t *testing.T) {
	t.Helper()
	err := os.RemoveAll("./docs")
	if err != nil {
		t.Fatalf("failed to clean up test CSV file: %v", err)
	}
}

func TestImportLocations(t *testing.T) {
	mockRepo := &MockDriverRepository{}
	service := NewDriverService(mockRepo)

	tests := []struct {
		name        string
		csvContent  string
		setupMock   func()
		expectedErr bool
	}{
		{
			name: "Successful Import",
			csvContent: `latitude,longitude
40.748817,-73.985428
34.052235,-118.243683`,
			setupMock: func() {
				mockRepo.On("SaveDrivers", mock.Anything, mock.Anything).Return(nil).Once()
			},
			expectedErr: false,
		},
		{
			name: "SaveDrivers Fails",
			csvContent: `latitude,longitude
40.748817,-73.985428
34.052235,-118.243683`,
			setupMock: func() {
				mockRepo.On("SaveDrivers", mock.Anything, mock.Anything).Return(errors.New("failed to save drivers")).Once()
			},
			expectedErr: true,
		},
		{
			name: "Invalid CSV Format",
			csvContent: `latitude,longitude
invalid,coordinates`,
			setupMock: func() {
				// No repository call expected
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create the temporary test CSV file
			_ = createTestCSVFile(t, tt.csvContent)
			defer removeTestCSVFile(t)

			// Set up mocks
			tt.setupMock()

			// Call the service with the file path
			err := service.ImportLocations(context.Background())
			if (err != nil) != tt.expectedErr {
				t.Errorf("expected error: %v, got: %v", tt.expectedErr, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
func TestFindNearestDriver(t *testing.T) {
	mockRepo := &MockDriverRepository{}
	service := NewDriverService(mockRepo)

	tests := []struct {
		name        string
		setupMock   func()
		latitude    float64
		longitude   float64
		radius      int
		expectedErr bool
	}{
		{
			name: "Successful Find",
			setupMock: func() {
				mockRepo.On("EnsureIndex", mock.Anything).Return(nil).Once()
				mockRepo.On("FindNearestDriver", mock.Anything, 40.748817, -73.985428, 5000).
					Return(&models.DriverWithDistance{
						Location: models.Location{
							Type:        "Point",
							Coordinates: []float64{-73.985428, 40.748817},
						},
						Distance: 100,
					}, nil).Once()
			},
			latitude:    40.748817,
			longitude:   -73.985428,
			radius:      5000,
			expectedErr: false,
		},
		{
			name: "No Drivers Found",
			setupMock: func() {
				mockRepo.On("EnsureIndex", mock.Anything).Return(nil).Once()
				mockRepo.On("FindNearestDriver", mock.Anything, 40.748817, -73.985428, 5000).
					Return(nil, repository.ErrDriverNotFound).Once()
			},
			latitude:    40.748817,
			longitude:   -73.985428,
			radius:      5000,
			expectedErr: true,
		},
		{
			name: "EnsureIndex Fails",
			setupMock: func() {
				mockRepo.On("EnsureIndex", mock.Anything).Return(errors.New("index creation failed")).Once()
			},
			latitude:    40.748817,
			longitude:   -73.985428,
			radius:      5000,
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			_, err := service.FindNearestDriver(context.Background(), tt.latitude, tt.longitude, tt.radius)
			if (err != nil) != tt.expectedErr {
				t.Errorf("expected error: %v, got: %v", tt.expectedErr, err)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}
