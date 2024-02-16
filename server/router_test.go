package server_test

import (
	"net/http"
	"testing"

	"github.com/belarte/metrix/model"
	"github.com/belarte/metrix/repository"
	"github.com/belarte/metrix/server"
	"github.com/playwright-community/playwright-go"
	"github.com/stretchr/testify/suite"
)

const address string = "127.0.0.1:12345"

var (
	initialMetrics = model.Metrics{
		{
			Title:       "Metric 1",
			Unit:        "unit",
			Description: "description",
		},
		{
			Title:       "Metric 2",
			Unit:        "unit",
			Description: "description",
		},
		{
			Title:       "Metric 3",
			Unit:        "unit",
			Description: "description",
		},
	}
	initialEntries = model.Entries{
		model.NewEntry(1, 5.0, "2018-01-01"),
		model.NewEntry(2, 2.1, "2018-01-11"),
		model.NewEntry(1, 1.0, "2018-01-15"),
	}
)

func populateDatabase(db *repository.Repository) error {
	for _, metric := range initialMetrics {
		_, err := db.UpsertMetric(metric)
		if err != nil {
			return err
		}
	}

	for _, entry := range initialEntries {
		_, err := db.UpsertEntry(entry)
		if err != nil {
			return err
		}
	}

	return nil
}

type RouterTestSuite struct {
	suite.Suite
	pw      *playwright.Playwright
	browser playwright.Browser
	context playwright.BrowserContext
	page    playwright.Page
	server  *server.Server
	db      *repository.Repository
}

func (s *RouterTestSuite) SetupSuite() {
	pw, err := playwright.Run()
	s.NoError(err)
	s.pw = pw

	browser, err := pw.Chromium.Launch()
	s.NoError(err)
	s.browser = browser

	context, err := browser.NewContext()
	s.NoError(err)
	context.SetDefaultTimeout(5000)
	s.context = context

	page, err := context.NewPage()
	s.NoError(err)
	s.page = page
}

func (s *RouterTestSuite) TearDownSuite() {
	err := s.context.Close()
	s.NoError(err)

	err = s.browser.Close()
	s.NoError(err)

	err = s.pw.Stop()
	s.NoError(err)
}

func (s *RouterTestSuite) SetupTest() {
	db, err := repository.New(":memory:")
	s.NoError(err)

	err = db.Migrate()
	s.NoError(err)

	err = populateDatabase(db)
	s.NoError(err)

	s.db = db
	s.server = server.New(
		server.WithRepository(db),
		server.WithAddress(address),
	)
	go func() {
		err := s.server.Start()
		s.ErrorIs(err, http.ErrServerClosed)
	}()
}

func (s *RouterTestSuite) TearDownTest() {
	err := s.server.Stop()
	s.NoError(err)

	err = s.db.Close()
	s.NoError(err)
}

func (s *RouterTestSuite) LoadPage() {
	_, err := s.page.Goto(address)
	s.NoError(err)

	err = s.page.Locator("text=Metrix 2024").WaitFor()
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

func (s *RouterTestSuite) TestAddMetric() {
	s.LoadPage()

	GoToManagePage(s.page, s.T()).
		Select("Create new metric").
		VerifyForm("", "", "").
		FillForm("new metric", "new unit", "new description").
		Create().
		Reload().
		Select("new metric").
		VerifyForm("new metric", "new unit", "new description")

	GoToEntryPage(s.page, s.T()).
		Select("new metric")
}

func (s *RouterTestSuite) TestUpdateMetric() {
	s.LoadPage()

	GoToManagePage(s.page, s.T()).
		Select("Metric 2").
		VerifyForm("Metric 2", "unit", "description").
		FillForm("Metric 2", "new unit", "new description").
		Update().
		Reload().
		Select("Metric 2").
		VerifyForm("Metric 2", "new unit", "new description")
}

func (s *RouterTestSuite) TestAddingEntryIsVisibleInReport() {
	s.LoadPage()

	GoToReportsPage(s.page, s.T()).
		Select("Metric 1").
		OpenEntriesList().
		VerifyEntriesCount(2)

	GoToEntryPage(s.page, s.T()).
		Select("Metric 1").
		AddEntry("2021-01-01", "7,0")

	GoToReportsPage(s.page, s.T()).
		Select("Metric 1").
		OpenEntriesList().
		VerifyEntriesCount(3).
		VerifyEntry("2021-01-01", "7.0")
}

func TestRouterTestSuite(t *testing.T) {
	suite.Run(t, new(RouterTestSuite))
}
