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

func (handler *ManageHandler) Manage(c echo.Context) error {
	metrics, err := handler.db.GetMetrics()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	page := views.ManagePage(metrics, model.Metric{}, submitButtonCreate, "")
	return render(c, page)
}

func (handler *ManageHandler) Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	err = handler.db.DeleteMetric(id)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	metrics, err := handler.db.GetMetrics()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	component := views.Manage(metrics, model.Metric{}, submitButtonCreate, "")
	return render(c, component)
}

func (handler *ManageHandler) Submit(c echo.Context) error {
	var metric model.Metric
	if err := c.Bind(&metric); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	upsertedMetric, err := handler.db.UpsertMetric(metric)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	additionnalMessage := "Metric updated!"
	buttonLabel := submitButtonUpdate
	if metric.ID == 0 {
		additionnalMessage = "Metric created!"
		buttonLabel = submitButtonCreate
	}

	metrics, err := handler.db.GetMetrics()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	page := views.Manage(metrics, upsertedMetric, buttonLabel, additionnalMessage)
	return render(c, page)
}

func (handler *ManageHandler) Select(c echo.Context) error {
	var metric model.Metric
	label := submitButtonCreate

	if val := c.FormValue("manage-select"); val != "create" {
		id, err := strconv.Atoi(val)
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

	page := views.Manage(metrics, metric, label, "")
	return render(c, page)
}
