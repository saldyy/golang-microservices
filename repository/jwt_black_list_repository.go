package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	JWT_BLACK_LIST_KEY = "jwt_black_list"
)

type JwtBlackListRepository struct {
	client *redis.Client
}

func (r *JwtBlackListRepository) Add(token string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	key := JWT_BLACK_LIST_KEY + ":" + token
	// This value should equal JWT token expiration time
  err := r.client.Set(ctx, key, 1, time.Hour).Err()

  if err != nil {
    return err
  }

  return nil
}

func (r *JwtBlackListRepository) Get(token string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	key := JWT_BLACK_LIST_KEY + ":" + token
	// This value should equal JWT token expiration time
  val, err := r.client.Get(ctx, key).Result()

  if err != nil {
    return "", err
  }

  return val, nil
}

func NewJwtBlackListRepository(client *redis.Client) *JwtBlackListRepository {

	return &JwtBlackListRepository{client: client}
}
