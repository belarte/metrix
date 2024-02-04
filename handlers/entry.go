package handlers

import (
	"net/http"
	"strconv"

	"github.com/belarte/metrix/database"
	"github.com/belarte/metrix/model"
	"github.com/belarte/metrix/views"
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

	page := views.EntryPage(metrics)
	return render(c, page)
}

func (h *EntryHandler) Select(c echo.Context) error {
	var metric model.Metric

	if val := c.FormValue("entry-select"); val != "add-entry" {
		id, err := strconv.Atoi(val)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		metric, err = h.db.GetMetric(id)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		form := views.EntryForm(metric)
		return render(c, form)
	}

	return c.String(http.StatusOK, "Please select a metric.")
}

func (h *EntryHandler) Submit(c echo.Context) error {
	var entry model.Entry
	if err := c.Bind(&entry); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	_, err := h.db.UpsertEntry(entry.MetricID, entry.Value, entry.Date)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	res := views.EntryCreated()
	return render(c, res)
}
