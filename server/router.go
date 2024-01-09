package server

import (
	"net/http"

	"github.com/belarte/metrix/database"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type templateParams struct {
	Metrics     []database.Metric
	Selected    database.Metric
	Content     string
	ButtonLabel string
}

func homeHandler(c echo.Context) error {
	return c.Render(http.StatusOK, "page", templateParams{
		Content: "home",
	})
}

func Run(addr string, db *database.InMemory) error {
	manageHandler := NewManageHandler(db)
	e := echo.New()

	e.Renderer = NewTemplateRenderer()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", homeHandler)
	e.GET("/manage", manageHandler.Manage)
	e.POST("/manage/submit", manageHandler.Submit)
	e.GET("/manage/select", manageHandler.Select)

	return e.Start(addr)
}
