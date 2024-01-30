package server_test

import (
	"testing"

	"github.com/playwright-community/playwright-go"
	"github.com/stretchr/testify/assert"
)

var pwa = playwright.NewPlaywrightAssertions()

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

func submit(page playwright.Page, t *testing.T, button, confirmation string) {
	err := page.GetByRole("button", playwright.PageGetByRoleOptions{
		Name: button,
	}).Click()
	assert.NoError(t, err)

	err = page.GetByText(confirmation).WaitFor()
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

func (p *ManagePage) VerifyForm(title, unit, description string) *ManagePage {
	err := pwa.Locator(p.page.GetByLabel("Title")).ToHaveValue(title)
	assert.NoError(p.t, err)

	err = pwa.Locator(p.page.GetByLabel("Unit")).ToHaveValue(unit)
	assert.NoError(p.t, err)

	err = pwa.Locator(p.page.GetByLabel("Description")).ToHaveValue(description)
	assert.NoError(p.t, err)

	return p
}

func (p *ManagePage) Reload() *ManagePage {
	p.page.Reload()
	return p
}

func (p *ManagePage) Create() *ManagePage {
	submit(p.page, p.t, "Create", "Metric created!")
	return p
}

func (p *ManagePage) Update() *ManagePage {
	submit(p.page, p.t, "Update", "Metric updated!")
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

func (p *EntryPage) Select(title string) *EntryPage {
	selectOption(p.page, p.t, title)
	return p
}

func (p *EntryPage) AddEntry(date, value string) *EntryPage {
	err := p.page.GetByText("Value in").PressSequentially(value)
	assert.NoError(p.t, err)

	err = p.page.GetByText("Date").Fill(date)
	assert.NoError(p.t, err)

	submit(p.page, p.t, "Submit", "Entry submitted successfully!")
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

func (p *ReportsPage) Select(title string) *ReportsPage {
	selectOption(p.page, p.t, title)
	return p
}

func (p *ReportsPage) OpenEntriesList() *ReportsPage {
	err := p.page.GetByText("Entries").Click()
	assert.NoError(p.t, err)
	return p
}

func (p *ReportsPage) VerifyEntriesCount(expectedCount int) *ReportsPage {
	count, err := p.page.Locator("table tbody tr").Count()
	assert.NoError(p.t, err)
	assert.Equal(p.t, expectedCount, count)

	return p
}

func (p *ReportsPage) VerifyEntry(date, value string) *ReportsPage {
    err := p.page.GetByText(date).WaitFor()
    assert.NoError(p.t, err)

    err = p.page.GetByText(value).WaitFor()
    assert.NoError(p.t, err)

	return p
}
