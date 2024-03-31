package repository

import (
	"context"

	"github.com/saldyy/golang-microservices/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var collectionName = "user"

type UserRepository struct {
	collection *mongo.Collection
}

func (r *UserRepository) FindByUsername(username string) (*model.User, error) {
	var result model.User
	err := r.collection.FindOne(context.TODO(), bson.M{"username": username}).Decode(&result)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func NewUserRepository(db *mongo.Database) *UserRepository {
  
	return &UserRepository{collection: db.Collection(collectionName)}
}
