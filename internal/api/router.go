package api

import (
	"net/http"

	"github.com/gorilla/mux"

	"bitaksi-go-driver/internal/api/handler"
	"bitaksi-go-driver/internal/config"
	"bitaksi-go-driver/internal/middleware"
)

// SetupRouter initializes the application router with all endpoints and middleware
func SetupRouter(driverService handler.DriverService, cfg *config.Config) *mux.Router {
	router := mux.NewRouter()

	// Initialize handlers
	driverHandler := handler.NewDriverHandler(driverService)

	// Public routes (e.g., health check)
	router.HandleFunc("/health", handler.HealthCheckHandler).Methods("GET")

	// Driver-related routes with API key middleware
	driverRouter := router.PathPrefix("/driver/api/v1").Subrouter()
	driverRouter.Use(middleware.APIKeyMiddleware(cfg))

	// Register driver endpoints
	driverRouter.HandleFunc("/import", driverHandler.ImportLocations).Methods(http.MethodPost)
	driverRouter.HandleFunc("/search", driverHandler.FindNearestDriver).Methods(http.MethodGet)

	return router
}
