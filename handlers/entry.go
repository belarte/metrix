package handlers

import (
	"net/http"
	"strconv"

	"github.com/belarte/metrix/model"
	"github.com/belarte/metrix/repository"
	"github.com/belarte/metrix/views"
	"github.com/labstack/echo/v4"
)

type EntryHandler struct {
	db *repository.Repository
}

func NewEntryHandler(db *repository.Repository) *EntryHandler {
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

		entries, err := h.db.GetSortedEntriesForMetric(id)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		form := views.EntriesTable(metric.ID, entries)
		return render(c, form)
	}

	return c.String(http.StatusOK, "Please select a metric.")
}

func (h *EntryHandler) Submit(c echo.Context) error {
	var entry model.Entry
	if err := c.Bind(&entry); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	upserted, err := h.db.UpsertEntry(entry)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	res := views.Row(upserted)
	return render(c, res)
}

func (h *EntryHandler) GetEntry(c echo.Context) error {
	var params model.Entry
	if err := c.Bind(&params); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	entry, err := h.db.GetEntry(params.MetricID, params.Date)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	res := views.Row(entry)
	return render(c, res)
}

func (h *EntryHandler) UpdateEntry(c echo.Context) error {
	var entry model.Entry
	if err := c.Bind(&entry); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	upserted, err := h.db.UpsertEntry(entry)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	res := views.Row(upserted)
	return render(c, res)
}

func (h *EntryHandler) GetEditableEntry(c echo.Context) error {
	var params model.Entry
	if err := c.Bind(&params); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	entry, err := h.db.GetEntry(params.MetricID, params.Date)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	res := views.EditableRow(entry)
	return render(c, res)
}

func (h *EntryHandler) Delete(c echo.Context) error {
	var params model.Entry
	if err := c.Bind(&params); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	if err := h.db.DeleteEntry(params.MetricID, params.Date); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}
