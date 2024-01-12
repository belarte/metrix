package server

import (
	"net/http"
	"strconv"

	"github.com/belarte/metrix/database"
	"github.com/labstack/echo/v4"
)

type EntryHandler struct {
	db *database.InMemory
}

func NewEntryHandler(db *database.InMemory) *EntryHandler {
	return &EntryHandler{
		db: db,
	}
}

func (h *EntryHandler) Entry(c echo.Context) error {
	metrics, err := h.db.GetMetrics()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.Render(http.StatusOK, "page", templateParams{
		Content: "entry",
		Metrics: metrics,
	})
}

func (h *EntryHandler) Select(c echo.Context) error {
	var metric database.Metric

	if val := c.FormValue("entry-select"); val != "add-entry" {
		id, err := strconv.Atoi(val)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		metric, err = h.db.GetMetric(id)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.Render(http.StatusOK, "entry-content", templateParams{
			Selected: metric,
		})
	}

	return c.String(http.StatusOK, "Please select a metric.")
}
