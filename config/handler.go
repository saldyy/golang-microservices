package config

import (
	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/saldyy/golang-microservices/repository"
)

type HandlerOptions struct {
	e      *echo.Echo
	repo   *repository.Repository
	logger *slog.Logger
}
