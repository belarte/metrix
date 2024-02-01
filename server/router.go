package server

import (
	"context"

	"github.com/belarte/metrix/database"
	"github.com/belarte/metrix/views"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type templateParams struct {
	Metrics           []database.Metric
	Selected          database.Metric
	Content           string
	ButtonLabel       string
	AdditionalMessage string
}

func homeHandler(c echo.Context) error {
	home := views.HomePage()
	return render(c, home)
}

type Server struct {
	db *database.InMemory
	e  *echo.Echo
}

func New(db *database.InMemory) *Server {
	manageHandler := NewManageHandler(db)
	entryHandler := NewEntryHandler(db)
	reportsHandler := NewReportsHandler(db)

	e := echo.New()
	e.Renderer = NewTemplateRenderer()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", homeHandler)
	e.GET("/manage", manageHandler.Manage)
	e.POST("/manage/submit", manageHandler.Submit)
	e.GET("/manage/select", manageHandler.Select)
	e.GET("/entry", entryHandler.Entry)
	e.GET("/entry/select", entryHandler.Select)
	e.POST("/entry/submit", entryHandler.Submit)
	e.GET("/reports", reportsHandler.Reports)
	e.GET("/reports/select", reportsHandler.Select)

	return &Server{
		db: db,
		e:  e,
	}
}

func (s *Server) Start(addr string) error {
	return s.e.Start(addr)
}

func (s *Server) Stop() error {
	return s.e.Shutdown(context.Background())
}
