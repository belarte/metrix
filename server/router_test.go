package server_test

import (
	"testing"

	"github.com/playwright-community/playwright-go"
	"github.com/stretchr/testify/assert"
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

func TestRouterLandsOnTheHomePage(t *testing.T) {
	pw, err := playwright.Run()
	assert.NoError(t, err)

	browser, err := pw.Chromium.Launch()
	assert.NoError(t, err)

	page, err := browser.NewPage()
	assert.NoError(t, err)

	_, err = page.Goto("localhost:8080")
	assert.NoError(t, err)

	err = page.Locator("text=Welcome").WaitFor()
	assert.NoError(t, err)

	err = browser.Close()
	assert.NoError(t, err)

	err = pw.Stop()
	assert.NoError(t, err)
}

func TestCanNavigateBetweenPages(t *testing.T) {
	pw, err := playwright.Run()
	assert.NoError(t, err)

	browser, err := pw.Chromium.Launch()
	assert.NoError(t, err)

	page, err := browser.NewPage()
	assert.NoError(t, err)

	_, err = page.Goto("localhost:8080")
	assert.NoError(t, err)

	_ = GoToHomePage(page, t)
	_ = GoToManagePage(page, t)
	_ = GoToEntryPage(page, t)
	_ = GoToReportsPage(page, t)

	err = browser.Close()
	assert.NoError(t, err)

	err = pw.Stop()
	assert.NoError(t, err)
}
