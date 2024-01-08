package server

import (
	"log"
	"net/http"
	"strconv"

	"github.com/belarte/metrix/database"
	"github.com/labstack/echo/v4"
)

type ManageHandler struct {
	db *database.InMemory
}

func NewManageHandler(db *database.InMemory) *ManageHandler {
	return &ManageHandler{
		db: db,
	}
}

func (handler *ManageHandler) Manage(c echo.Context) error {
	metrics, err := handler.db.GetMetrics()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.Render(http.StatusOK, "page", templateParams{
		Metrics:  metrics,
		Selected: database.Metric{},
		Content:  "manage",
	})
}

func (handler *ManageHandler) Submit(c echo.Context) error {
	var metric database.Metric
	if err := c.Bind(&metric); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	metric, err := handler.db.AddMetric(metric)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	metrics, err := handler.db.GetMetrics()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.Render(http.StatusOK, "manage", templateParams{
		Metrics:  metrics,
		Selected: metric,
	})
}

func (handler *ManageHandler) Select(c echo.Context) error {
	val := c.FormValue("manage-select")

	id, err := strconv.Atoi(c.FormValue("manage-select"))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	log.Printf("val: %s - id: %d", val, id)

	metric, err := handler.db.GetMetric(id)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	metrics, err := handler.db.GetMetrics()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.Render(http.StatusOK, "manage", templateParams{
		Metrics:  metrics,
		Selected: metric,
	})
}
