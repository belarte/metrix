package handlers

import (
	"net/http"
	"strconv"

	"github.com/belarte/metrix/model"
	"github.com/belarte/metrix/repository"
	"github.com/belarte/metrix/views"
	"github.com/labstack/echo/v4"
)

type ManageHandler struct {
	db *repository.Repository
}

func NewManageHandler(db *repository.Repository) *ManageHandler {
	return &ManageHandler{
		db: db,
	}
}

var (
	submitButtonCreate = "Create"
	submitButtonUpdate = "Update"
)

func (h *ManageHandler) Manage(c echo.Context) error {
	metrics, err := h.db.GetMetrics()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	page := views.ManagePage(metrics, model.Metric{}, submitButtonCreate, "")
	return render(c, page)
}

func (h *ManageHandler) Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	err = h.db.DeleteMetric(id)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return h.renderManageComponent(c, model.Metric{}, submitButtonCreate, "")
}

func (h *ManageHandler) Submit(c echo.Context) error {
	var metric model.Metric
	if err := c.Bind(&metric); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	upsertedMetric, err := h.db.UpsertMetric(metric)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	additionnalMessage := "Metric updated!"
	buttonLabel := submitButtonUpdate
	if metric.ID == 0 {
		additionnalMessage = "Metric created!"
		buttonLabel = submitButtonCreate
	}

	return h.renderManageComponent(c, upsertedMetric, buttonLabel, additionnalMessage)
}

func (h *ManageHandler) Select(c echo.Context) error {
	var metric model.Metric
	label := submitButtonCreate

	if val := c.FormValue("manage-select"); val != "create" {
		id, err := strconv.Atoi(val)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		metric, err = h.db.GetMetric(id)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		label = submitButtonUpdate
	}

	return h.renderManageComponent(c, metric, label, "")
}

func (h *ManageHandler) renderManageComponent(c echo.Context, m model.Metric, label, message string) error {
	metrics, err := h.db.GetMetrics()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	component := views.Manage(metrics, m, label, message)
	return render(c, component)
}
