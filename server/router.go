package server

import (
	"context"

	"github.com/belarte/metrix/handlers"
	"github.com/belarte/metrix/repository"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	db      *repository.Repository
	e       *echo.Echo
	address string
}

type Option func(*Server)

func WithAddress(addr string) Option {
	return func(s *Server) {
		s.address = addr
	}
}

func WithRepository(db *repository.Repository) Option {
	return func(s *Server) {
		s.db = db
	}
}

func newServer(db *repository.Repository) *echo.Echo {
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

	return e
}

func New(options ...Option) *Server {
	s := &Server{}
	for _, option := range options {
		option(s)
	}

	s.e = newServer(s.db)
	return s
}

func (s *Server) Start() error {
	return s.e.Start(s.address)
}

func (s *Server) Stop() error {
	return s.e.Shutdown(context.Background())
}
