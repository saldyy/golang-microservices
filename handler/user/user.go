package user

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/saldyy/golang-microservices/repository"
)

type UserHandler struct {
	userRepo *repository.UserRepository
	logger   *slog.Logger
}

func (u *UserHandler) MeHandler(c echo.Context) error {
	return c.String(http.StatusOK, "Me")
}

func RegisterUserHandlers(e *echo.Echo, repo *repository.Repository, logger *slog.Logger) *UserHandler {
	userHandler := UserHandler{userRepo: repo.UserRepository, logger: logger}

	e.GET("/me", userHandler.MeHandler)

	return &userHandler
}
