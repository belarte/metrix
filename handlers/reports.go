package handlers

import (
	"net/http"
	"strconv"

	"github.com/belarte/metrix/database"
	"github.com/belarte/metrix/diagram"
	"github.com/belarte/metrix/views"
	"github.com/labstack/echo/v4"
)

type ReportsHandler struct {
	db *database.InMemory
}

func NewReportsHandler(db *database.InMemory) *ReportsHandler {
	return &ReportsHandler{
		db: db,
	}
}

func (h *ReportsHandler) Reports(c echo.Context) error {
	metrics, err := h.db.GetMetrics()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	reports := views.ReportsPage(metrics)
	return render(c, reports)
}

func (h *ReportsHandler) Select(c echo.Context) error {
	val := c.FormValue("reports-select")
	if val == "reports-select" {
		return c.String(http.StatusOK, "Please select a metric.")
	}

	id, err := strconv.Atoi(val)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	metric, err := h.db.GetMetric(id)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	entries, err := h.db.GetSortedEntriesForMetric(id)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	s := diagram.DataToGraph(metric, entries)
	reports := views.Reports(entries, s)

	return render(c, reports)
}
