package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Instance *DatabaseInstance

type DatabaseInstance struct {
	DB *mongo.Database
}

func InitMongoClient() *DatabaseInstance {
	connectionUri := os.Getenv("MONGODB_URI")
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(connectionUri))

	if err != nil {
		log.Fatal("Cannot connect to Mongodb\n")
		panic(err)
	}
	fmt.Printf("Connected to MongoDB\n")
	return &DatabaseInstance{DB: client.Database("test")}
}

func (di *DatabaseInstance) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := di.DB.Client().Ping(ctx, nil)
	if err != nil {
		log.Fatalf(fmt.Sprintf("db down: %v", err))
	}

	return map[string]string{
		"message": "ok",
	}
}
