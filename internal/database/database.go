package database

import (
	"context"
	"log/slog"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Instance *DatabaseInstance

type DatabaseInstance struct {
	DB *mongo.Database
  logger *slog.Logger
}

func InitMongoClient(logger *slog.Logger) *DatabaseInstance {
	connectionUri := os.Getenv("MONGODB_URI")
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(connectionUri))

	if err != nil {
    logger.Error("Error connecting to MongoDB", err)
		panic(err)
	}
  logger.Info("Connected to MongoDB")
	return &DatabaseInstance{DB: client.Database("test"), logger: logger}
}

func (di *DatabaseInstance) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := di.DB.Client().Ping(ctx, nil)
	if err != nil {
    di.logger.Error("Error connecting to MongoDB", err)
		return map[string]string{
			"message": "not_ok",
		}
	}

	return map[string]string{
		"message": "ok",
	}
}
