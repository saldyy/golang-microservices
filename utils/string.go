package utils

import (
	"golang.org/x/crypto/bcrypt"
)

const (
	MinCost     int = 4
	MaxCost     int = 31
	DefaultCost int = 10
)

func BcryptHash(str string) (string, error) {
  hash, err := bcrypt.GenerateFromPassword([]byte(str), DefaultCost)
  if err != nil {
    return "", err
  }
  return string(hash), nil
}
