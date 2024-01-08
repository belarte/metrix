package server

import (
	"html/template"
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
	t, err := template.ParseFiles("server/templates/main.tmpl", "server/templates/home.tmpl")
	if err != nil {
		log.Fatal(err)
	}
	return t.Execute(c.Response().Writer, templateParams{Content: "content"})
}

func manageHandler(db *database.InMemory) func(echo.Context) error {
	return func(c echo.Context) error {
		t, err := template.ParseFiles("server/templates/main.tmpl", "server/templates/manage.tmpl")
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		metrics, err := db.GetMetrics()
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return t.Execute(c.Response().Writer, templateParams{
			Metrics:  metrics,
			Selected: database.Metric{},
			Content:  "content",
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

		t, err := template.ParseFiles("server/templates/manage.tmpl")
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		metrics, err := db.GetMetrics()
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return t.ExecuteTemplate(c.Response().Writer, "content", templateParams{
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

		t, err := template.ParseFiles("server/templates/manage.tmpl")
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return t.ExecuteTemplate(c.Response().Writer, "content", templateParams{
			Metrics:  metrics,
			Selected: metric,
		})
	}
}

func Run(addr string, db *database.InMemory) error {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.GET("/", homeHandler)
	e.GET("/manage", manageHandler(db))
	e.POST("/manage/click", clickHandler(db))
	e.GET("/manage/select", selectHandler(db))
	return e.Start(addr)
}
