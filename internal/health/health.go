package health

import (
	"net/http"

	"github.com/labstack/echo/v4"
)


func CheckHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"status": "up"})

}
