package repository

import "go.mongodb.org/mongo-driver/mongo"

type Repository struct{
  UserRepository *UserRepository
}

func Init(db *mongo.Database) *Repository {
	return &Repository{
    UserRepository: NewUserRepository(db),
  }
}
