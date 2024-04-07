package repository

import (
	"context"

	"github.com/saldyy/golang-microservices/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (r *UserRepository) InsertUser(user *model.User) (*model.User, error) {
	var result *mongo.InsertOneResult
	// result, err := r.collection.InsertOne(context.TODO(), bson.D{{"username", user.Username}, {"password", user.Password}})
  result, err := r.collection.InsertOne(context.TODO(), &user)

	if err != nil {
		return nil, err
	}

	user.Id = result.InsertedID.(primitive.ObjectID).Hex()

	return user, nil
}

func NewUserRepository(db *mongo.Database) *UserRepository {

	return &UserRepository{collection: db.Collection(collectionName)}
}
