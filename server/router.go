package server

import (
	"context"

	"github.com/belarte/metrix/database"
	"github.com/belarte/metrix/handlers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	db *database.InMemory
	e  *echo.Echo
}

func New(db *database.InMemory) *Server {
	manageHandler := handlers.NewManageHandler(db)
	entryHandler := handlers.NewEntryHandler(db)
	reportsHandler := handlers.NewReportsHandler(db)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", handlers.HomeHandler)
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
