package health

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/saldyy/golang-microservices/internal/database"
)

type HealthHandler struct {
	db *database.Db
}

func New(db *database.Db) *HealthHandler {
	return &HealthHandler{db: db}
}

func (h *HealthHandler) CheckHandler(c echo.Context) error {
	dbStat := h.checkDb(c)

	return c.JSON(http.StatusOK, map[string]string{"status": "up", "databaseStatus": dbStat})
}

func (h *HealthHandler) checkDb(c echo.Context) string {
	if h.db == nil {
		return "down"
	}

	dbStat, err := h.db.Check()
	if err != nil {
		c.Logger().Errorf("Database down: %v", err)
	}

	return dbStat
}
