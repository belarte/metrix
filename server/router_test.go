package server_test

import (
	"net/http"
	"testing"

	"github.com/belarte/metrix/database"
	"github.com/belarte/metrix/server"
	"github.com/playwright-community/playwright-go"
	"github.com/stretchr/testify/suite"
)

var address = "127.0.0.1:12345"

type RouterTestSuite struct {
	suite.Suite
	pw      *playwright.Playwright
	browser playwright.Browser
	page    playwright.Page
	server  *server.Server
}

func (s *RouterTestSuite) SetupSuite() {
	pw, err := playwright.Run()
	s.NoError(err)
	s.pw = pw

	browser, err := pw.Chromium.Launch()
	s.NoError(err)
	s.browser = browser

	page, err := browser.NewPage()
	s.NoError(err)
	s.page = page
}

func (s *RouterTestSuite) TearDownSuite() {
	err := s.browser.Close()
	s.NoError(err)

	err = s.pw.Stop()
	s.NoError(err)
}

func (s *RouterTestSuite) SetupTest() {
	s.server = server.New(database.NewInMemory())
	go func() {
		err := s.server.Start(address)
		s.ErrorIs(err, http.ErrServerClosed)
	}()
}

func (s *RouterTestSuite) TearDownTest() {
	err := s.server.Stop()
	s.NoError(err)
}

func (s *RouterTestSuite) LoadPage() {
	_, err := s.page.Goto(address)
	s.NoError(err)
}

func (s *RouterTestSuite) TestRouterLandsOnTheHomePage() {
	s.LoadPage()

	err := s.page.Locator("text=Welcome").WaitFor()
	s.NoError(err)
}

func (s *RouterTestSuite) TestRouterNavigatesBetweenPages() {
	s.LoadPage()

	_ = GoToReportsPage(s.page, s.T())
	_ = GoToEntryPage(s.page, s.T())
	_ = GoToManagePage(s.page, s.T())
	_ = GoToHomePage(s.page, s.T())
}

func TestRouterTestSuite(t *testing.T) {
	suite.Run(t, new(RouterTestSuite))
}
