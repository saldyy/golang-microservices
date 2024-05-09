package repository

import (
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct{
  // Redis
  JwtBlackListRepository *JwtBlackListRepository

  // Mongo
  UserRepository *UserRepository
  
}

func Init(db *mongo.Database, redisClient *redis.Client) *Repository {
	return &Repository{
    UserRepository: NewUserRepository(db),
    JwtBlackListRepository: NewJwtBlackListRepository(redisClient),
  }
}


