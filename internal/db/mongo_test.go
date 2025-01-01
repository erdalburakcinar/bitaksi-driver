package db

import (
	"bitaksi-go-driver/internal/config"
	"errors"
	"testing"
)

// MockClient simulates a MongoDB client for testing purposes
type MockClient struct {
	pingErr error
}

// ConnectMongoMocked mocks the ConnectMongo function
func ConnectMongoMocked(cfg *config.Config, client *MockClient) (*MockClient, error) {
	if client.pingErr != nil {
		return nil, client.pingErr
	}
	return client, nil
}

func TestConnectMongo(t *testing.T) {
	tests := []struct {
		name        string
		cfg         *config.Config
		mockClient  *MockClient
		expectError bool
	}{
		{
			name: "Successful Connection",
			cfg: &config.Config{
				MongoDB: struct {
					Username   string `mapstructure:"username"`
					Password   string `mapstructure:"password"`
					Host       string `mapstructure:"host"`
					Port       int    `mapstructure:"port"`
					Database   string `mapstructure:"database"`
					Collection string `mapstructure:"collection"`
				}{
					Username: "test-user",
					Password: "test-pass",
					Host:     "localhost",
					Port:     27017,
				},
			},
			mockClient:  &MockClient{pingErr: nil},
			expectError: false,
		},
		{
			name: "Ping Failure",
			cfg: &config.Config{
				MongoDB: struct {
					Username   string `mapstructure:"username"`
					Password   string `mapstructure:"password"`
					Host       string `mapstructure:"host"`
					Port       int    `mapstructure:"port"`
					Database   string `mapstructure:"database"`
					Collection string `mapstructure:"collection"`
				}{
					Username: "test-user",
					Password: "test-pass",
					Host:     "localhost",
					Port:     27017,
				},
			},
			mockClient:  &MockClient{pingErr: errors.New("ping failed")},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := ConnectMongoMocked(tt.cfg, tt.mockClient)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected an error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}

				if client == nil {
					t.Errorf("expected client but got nil")
				}
			}
		})
	}
}
