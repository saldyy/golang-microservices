package user

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/saldyy/golang-microservices/model"
	"github.com/saldyy/golang-microservices/repository"
	"github.com/saldyy/golang-microservices/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserHandler struct {
	userRepo *repository.UserRepository
	logger   *slog.Logger
}

type LoginCredential struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (u *UserHandler) LoginHandler(c echo.Context) error {
	credential := new(LoginCredential)

	if err := c.Bind(credential); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid credential"})
	}

	user, err := u.userRepo.FindByUsername(credential.Username)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			u.logger.Debug("User not found", err)
			return c.JSON(http.StatusForbidden, map[string]string{"message": "User not found"})
		}

		return c.String(http.StatusOK, "Error")
	}

  isPasswordMatch := utils.IsMatchBcryptHash(user.Password, credential.Password)

  if !isPasswordMatch {
    return c.JSON(http.StatusForbidden, map[string]string{"message": "Invalid password"})
  }

	return c.String(http.StatusOK, user.Username)
}

func (u *UserHandler) RegisterHandler(c echo.Context) error {
	credential := new(LoginCredential)

  if err := c.Bind(credential); err != nil {
    return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid payload"})
  }

	existedUser, err := u.userRepo.FindByUsername(credential.Username)

	if existedUser != nil {
		u.logger.Error("User already existed", err)
		return c.JSON(http.StatusConflict, map[string]string{"message": "User already existed"})
	}

	hashedPassword, err := utils.BcryptHash(credential.Password)

	if err != nil {
		u.logger.Error("Error hashing password", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Internal server error"})
	}

	user := &model.User{Username: credential.Username, Password: hashedPassword}
	_, err = u.userRepo.InsertUser(user)

	if err != nil {
		u.logger.Error("Error inserting user", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Internal server error"})
	}

	return c.String(http.StatusOK, "user.Username")
}

func RegisterUserHandlers(e *echo.Echo, repo *repository.Repository, logger *slog.Logger) *UserHandler {
	userHandler := UserHandler{userRepo: repo.UserRepository, logger: logger}

	e.POST("/login", userHandler.LoginHandler)
	e.POST("/register", userHandler.RegisterHandler)

	return &userHandler
}
