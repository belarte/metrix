package server

import (
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

var (
	submitButtonCreate = "Create"
	submitButtonUpdate = "Update"
)

func (handler *ManageHandler) Manage(c echo.Context) error {
	metrics, err := handler.db.GetMetrics()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.Render(http.StatusOK, "page", templateParams{
		Metrics:     metrics,
		Selected:    database.Metric{},
		Content:     "manage",
		ButtonLabel: submitButtonCreate,
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
		Metrics:     metrics,
		Selected:    metric,
		ButtonLabel: submitButtonUpdate,
	})
}

func (handler *ManageHandler) Select(c echo.Context) error {
	var metric database.Metric
	label := submitButtonCreate

	if val := c.FormValue("manage-select"); val != "create" {
		id, err := strconv.Atoi(c.FormValue("manage-select"))
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		metric, err = handler.db.GetMetric(id)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		label = submitButtonUpdate
	}

	metrics, err := handler.db.GetMetrics()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.Render(http.StatusOK, "manage", templateParams{
		Metrics:     metrics,
		Selected:    metric,
		ButtonLabel: label,
	})
}
