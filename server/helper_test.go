package server_test

import (
	"testing"

	"github.com/playwright-community/playwright-go"
	"github.com/stretchr/testify/assert"
)

func goToPage(page playwright.Page, t *testing.T, button, expectedTitle string) {
	err := page.GetByText(button).Click()
	assert.NoError(t, err)

	err = page.GetByText(expectedTitle).WaitFor()
	assert.NoError(t, err)
}

func selectOption(page playwright.Page, t *testing.T, label string) {
	values, err := page.GetByRole("combobox").SelectOption(playwright.SelectOptionValues{
		Labels: playwright.StringSlice(label),
	})
	assert.NoError(t, err)
	assert.Len(t, values, 1)
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

func (p *ManagePage) Select(label string) *ManagePage {
	selectOption(p.page, p.t, label)
	return p
}

func (p *ManagePage) FillForm(title, unit, description string) *ManagePage {
	err := p.page.GetByText("Title").Fill(title)
	assert.NoError(p.t, err)

	err = p.page.GetByText("Unit").Fill(unit)
	assert.NoError(p.t, err)

	err = p.page.GetByText("Description").Fill(description)
	assert.NoError(p.t, err)

	return p
}

func (p *ManagePage) Click(name string) *ManagePage {
	err := p.page.GetByRole("button", playwright.PageGetByRoleOptions{
		Name: name,
	}).Click()
	assert.NoError(p.t, err)

	return p
}

type EntryPage struct {
	page playwright.Page
	t    *testing.T
}

func GoToEntryPage(page playwright.Page, t *testing.T) *EntryPage {
	goToPage(page, t, "Entry", "Add an entry")
	return &EntryPage{page: page, t: t}
}

func (p *EntryPage) SelectMetric(title string) *EntryPage {
	selectOption(p.page, p.t, title)
	return p
}

type ReportsPage struct {
	page playwright.Page
	t    *testing.T
}

func GoToReportsPage(page playwright.Page, t *testing.T) *ReportsPage {
	goToPage(page, t, "Reports", "Consult reports")
	return &ReportsPage{page: page, t: t}
}
