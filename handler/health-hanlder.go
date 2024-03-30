package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/saldyy/golang-microservices/config/database"
)

func HealthCheckHandler(c echo.Context) error {
	healthStatus := make(map[string]map[string]string)

	healthStatus["database"] = database.Instance.Health()
	healthStatus["server"] = map[string]string{"message": "ok"}
	return c.JSON(http.StatusOK, healthStatus)
}
