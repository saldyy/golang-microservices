package handler

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/saldyy/golang-microservices/config/database"
)

type HealthHandler struct {
  e *echo.Echo
  logger *log.Logger
  db *database.DatabaseInstance
}


func HealthCheckHandler(c echo.Context) error {
	healthStatus := make(map[string]map[string]string)

	healthStatus["database"] = database.Instance.Health()
	healthStatus["server"] = map[string]string{"message": "ok"}
	return c.JSON(http.StatusOK, healthStatus)
}

func (h *HealthHandler) RegisterHealthCheckHandler(e *echo.Echo, logger *log.Logger, db *database.DatabaseInstance) {
  
  e.GET("/health", HealthCheckHandler)
}
