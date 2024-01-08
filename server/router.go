package server

import (
	"log"
	"net/http"
	"strconv"

	"github.com/belarte/metrix/database"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type templateParams struct {
	Metrics  []database.Metric
	Selected database.Metric
	Content  string
}

func homeHandler(c echo.Context) error {
	return c.Render(http.StatusOK, "page", templateParams{
		Content: "home",
	})
}

func manageHandler(db *database.InMemory) func(echo.Context) error {
	return func(c echo.Context) error {
		metrics, err := db.GetMetrics()
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.Render(http.StatusOK, "page", templateParams{
			Metrics:  metrics,
			Selected: database.Metric{},
			Content:  "manage",
		})
	}
}

func clickHandler(db *database.InMemory) func(echo.Context) error {
	return func(c echo.Context) error {
		var metric database.Metric
		if err := c.Bind(&metric); err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		metric, err := db.AddMetric(metric)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		metrics, err := db.GetMetrics()
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.Render(http.StatusOK, "manage", templateParams{
			Metrics:  metrics,
			Selected: metric,
		})
	}
}

func selectHandler(db *database.InMemory) func(echo.Context) error {
	return func(c echo.Context) error {
		val := c.FormValue("manage-select")

		id, err := strconv.Atoi(c.FormValue("manage-select"))
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		log.Printf("val: %s - id: %d", val, id)

		metric, err := db.GetMetric(id)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		metrics, err := db.GetMetrics()
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.Render(http.StatusOK, "manage", templateParams{
			Metrics:  metrics,
			Selected: metric,
		})
	}
}

func Run(addr string, db *database.InMemory) error {
	e := echo.New()
	e.Renderer = NewTemplateRenderer()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.GET("/", homeHandler)
	e.GET("/manage", manageHandler(db))
	e.POST("/manage/click", clickHandler(db))
	e.GET("/manage/select", selectHandler(db))
	return e.Start(addr)
}
