package handler

import (
	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/saldyy/golang-microservices/handler/user"
	"github.com/saldyy/golang-microservices/repository"
)

func RegisterRoutes(e *echo.Echo, repo *repository.Repository, logger *slog.Logger) {
	user.RegisterAuthHandlers(e, repo, logger)
	user.RegisterUserHandlers(e, repo, logger)
}
