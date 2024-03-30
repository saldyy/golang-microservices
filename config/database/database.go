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
	mongoClient *mongo.Client
}

func InitMongoClient() *DatabaseInstance {
	connectionUri := os.Getenv("MONGODB_URI")
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(connectionUri))

	if err != nil {
		log.Fatal("Cannot connect to Mongodb")
		panic(err)
	}
	fmt.Printf("Connected to MongoDB")
	return &DatabaseInstance{mongoClient: client}
}

func (s DatabaseInstance) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := s.mongoClient.Ping(ctx, nil)
	if err != nil {
		log.Fatalf(fmt.Sprintf("db down: %v", err))
	}

	return map[string]string{
		"message": "ok",
	}
}
