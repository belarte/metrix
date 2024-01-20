package server_test

import (
	"testing"

	"github.com/playwright-community/playwright-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type HomePage struct {
	page playwright.Page
	t    *testing.T
}

func goToPage(page playwright.Page, t *testing.T, button, expectedTitle string) {
	err := page.Locator("text=" + button).Click()
	assert.NoError(t, err)

	err = page.Locator("text=" + expectedTitle).WaitFor()
	assert.NoError(t, err)
}

func GoToHomePage(page playwright.Page, t *testing.T) *HomePage {
	goToPage(page, t, "Home", "Welcome")
	return &HomePage{page: page, t: t}
}

type ManagePage struct {
	page playwright.Page
	t    *testing.T
}

func GoToManagePage(page playwright.Page, t *testing.T) *ManagePage {
	goToPage(page, t, "Manage", "Manage metrics")
	return &ManagePage{page: page, t: t}
}

type EntryPage struct {
	page playwright.Page
	t    *testing.T
}

func GoToEntryPage(page playwright.Page, t *testing.T) *EntryPage {
	goToPage(page, t, "Entry", "Add an entry")
	return &EntryPage{page: page, t: t}
}

type ReportsPage struct {
	page playwright.Page
	t    *testing.T
}

func GoToReportsPage(page playwright.Page, t *testing.T) *ReportsPage {
	goToPage(page, t, "Reports", "Consult reports")
	return &ReportsPage{page: page, t: t}
}

type RouterTestSuite struct {
	suite.Suite
	pw      *playwright.Playwright
	browser playwright.Browser
	page    playwright.Page
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

func (s *RouterTestSuite) TestRouterLandsOnTheHomePage() {
	_, err := s.page.Goto("localhost:8080")
	s.NoError(err)

	err = s.page.Locator("text=Welcome").WaitFor()
	s.NoError(err)
}

func (s *RouterTestSuite) TestRouterNavigatesBetweenPages() {
	_ = GoToHomePage(s.page, s.T())
	_ = GoToManagePage(s.page, s.T())
	_ = GoToEntryPage(s.page, s.T())
	_ = GoToReportsPage(s.page, s.T())
}

func TestRouterTestSuite(t *testing.T) {
	suite.Run(t, new(RouterTestSuite))
}
