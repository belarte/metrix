package server_test

import (
	"testing"

	"github.com/playwright-community/playwright-go"
	"github.com/stretchr/testify/assert"
)

func goToPage(page playwright.Page, t *testing.T, button, expectedTitle string) {
	err := page.Locator("text=" + button).Click()
	assert.NoError(t, err)

	err = page.Locator("text=" + expectedTitle).WaitFor()
	assert.NoError(t, err)
}

type HomePage struct {
	page playwright.Page
	t    *testing.T
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
