package database

import (
	"context"
	"log"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload" // Automatically load environment variables from .env file
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Database struct {
	Client *mongo.Database
}

var (
	dbURI  = os.Getenv("DATABASE_URI")
	dbName = os.Getenv("DATABASE_NAME")
)

func NewDatabase() *Database {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(dbURI)

	client, err := mongo.Connect(clientOptions)
	if err != nil {
		log.Fatalf("Failed to create database client: %v", err)
	}

	if err = client.Ping(ctx, nil); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	return &Database{
		Client: client.Database(dbName),
	}
}
