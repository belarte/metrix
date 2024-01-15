package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/belarte/metrix/database"
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

	return c.Render(http.StatusOK, "page", templateParams{
		Content: "reports",
		Metrics: metrics,
	})
}

func (h *ReportsHandler) Select(c echo.Context) error {
	val := c.FormValue("reports-select")
	if val == "select-reports" {
        return c.String(http.StatusOK, "Please select a metric.")
	}

    id, err := strconv.Atoi(val)
    if err != nil {
        return c.String(http.StatusBadRequest, err.Error())
    }

    entries, err := h.db.GetEntriesForMetric(id)
    if err != nil {
        return c.String(http.StatusInternalServerError, err.Error())
    }

    s := fmt.Sprintf("entries: %+v\n", entries)
    return c.String(http.StatusOK, s)
}
