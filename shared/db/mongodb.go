package db

import (
	"context"
	"fmt"
	"log"
	"ride-sharing/shared/env"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	TripsCollection     = "trips"
	RideFaresCollection = "ride_fares"
)

type MongoConfig struct {
	URI      string
	Database string
}

func NewMongoDefaultConfig() *MongoConfig {
	return &MongoConfig{
		URI:      env.GetString("MONGODB_URI", "mongodb://root:example@mongo:27017/"),
		Database: "ride-sharing",
	}
}

func NewMongoClient(ctx context.Context, config *MongoConfig) (*mongo.Client, error) {
	if config.URI == "" {
		return nil, fmt.Errorf("mongodb URI is required")
	}
	if config.Database == "" {
		return nil, fmt.Errorf("mongodb database is required")
	}

	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.URI))
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	log.Printf("Successfully connected to MongoDB at %s", config.URI)
	return client, nil
}

func GetDatabase(client *mongo.Client, config *MongoConfig) *mongo.Database {
	return client.Database(config.Database)
}
