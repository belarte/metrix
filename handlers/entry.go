package handlers

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

		return c.Render(http.StatusOK, "entry-form", templateParams{
			Selected: metric,
		})
	}

	return c.String(http.StatusOK, "Please select a metric.")
}

func (h *EntryHandler) Submit(c echo.Context) error {
	var entry database.Entry
	if err := c.Bind(&entry); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	_, err := h.db.UpsertEntry(entry.MetricID, entry.Value, entry.Date)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.Render(http.StatusOK, "entry-created", nil)
}
