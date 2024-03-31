package user

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/saldyy/golang-microservices/repository"
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
  fmt.Printf("User: %v\n", user)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			u.logger.Debug("User not found", err)
			return c.String(http.StatusNotFound, "Notfound")
		}

		return c.String(http.StatusOK, "Error")
	}

	return c.String(http.StatusOK, "Signin")
}

func RegisterUserHandlers(e *echo.Echo, repo *repository.Repository, logger *slog.Logger) *UserHandler {
	userHandler := UserHandler{userRepo: repo.UserRepository, logger: logger}

	e.POST("/login", userHandler.LoginHandler)

	return &userHandler
}
