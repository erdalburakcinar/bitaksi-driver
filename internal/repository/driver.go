package repository

import (
	"bitaksi-go-driver/internal/models"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ErrDriverNotFound is returned when no driver is found within the specified radius
var ErrDriverNotFound = errors.New("no drivers found within the specified radius")

type DriverRepository struct {
	collection *mongo.Collection
}

func NewDriverRepository(db *mongo.Database, collectionName string) DriverRepository {
	return DriverRepository{collection: db.Collection(collectionName)}
}

func (r *DriverRepository) SaveDrivers(ctx context.Context, locations []models.DriverWithDistance) error {
	var drivers []interface{}

	for _, location := range locations {
		drivers = append(drivers, models.DriverWithDistance{
			ID:       primitive.NewObjectID(),
			Location: location.Location,
		})
	}

	_, err := r.collection.InsertMany(ctx, drivers)
	return err
}

func (r *DriverRepository) FindNearestDriver(ctx context.Context, latitude, longitude float64, maxDistance int) (*models.DriverWithDistance, error) {
	var drivers []models.DriverWithDistance

	// Define the geoNear aggregation pipeline
	aggregate := mongo.Pipeline{
		{
			{"$geoNear", bson.M{
				"near": bson.M{
					"type":        "Point",
					"coordinates": []float64{longitude, latitude},
				},
				"distanceField": "distance",  // Add the calculated distance
				"maxDistance":   maxDistance, // Maximum distance in meters
				"spherical":     true,        // Use spherical calculations
			}},
		},
		{{
			"$limit", 1, // Limit results to the nearest driver
		}},
	}

	// Execute the aggregation pipeline
	cursor, err := r.collection.Aggregate(ctx, aggregate)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Decode the result into the struct
	if err = cursor.All(ctx, &drivers); err != nil {
		return nil, err
	}

	// Check if a driver was found
	if len(drivers) == 0 {
		return nil, ErrDriverNotFound
	}

	// Return the nearest driver
	return &drivers[0], nil
}

// EnsureIndex ensures that the collection has a 2dsphere index on the location field.
func (r *DriverRepository) EnsureIndex(ctx context.Context) error {
	// Check if the 2dsphere index already exists
	indexes, err := r.collection.Indexes().List(ctx)
	if err != nil {
		return err
	}

	hasGeoIndex := false
	for indexes.Next(ctx) {
		var index bson.M
		if err := indexes.Decode(&index); err != nil {
			return err
		}

		if key, ok := index["key"].(bson.M); ok {
			if key["location"] == "2dsphere" {
				hasGeoIndex = true
				break
			}
		}
	}

	// Create the 2dsphere index if it does not exist
	if !hasGeoIndex {
		index := mongo.IndexModel{
			Keys:    bson.M{"location": "2dsphere"},
			Options: options.Index().SetName("location_2dsphere"),
		}

		_, err = r.collection.Indexes().CreateOne(ctx, index)
		if err != nil {
			return err
		}
	}
	return nil
}
