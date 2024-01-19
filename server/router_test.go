package server_test

import (
	"testing"

	"github.com/playwright-community/playwright-go"
	"github.com/stretchr/testify/assert"
)

func TestRouter(t *testing.T) {
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

	err = page.Locator("text=Manage").Click()
	assert.NoError(t, err)

	err = page.Locator("text=Manage metrics").WaitFor()
	assert.NoError(t, err)

	err = page.Locator("text=Entry").Click()
	assert.NoError(t, err)

	err = page.Locator("text=Add an entry").WaitFor()
	assert.NoError(t, err)

	err = page.Locator("text=Reports").Click()
	assert.NoError(t, err)

	err = page.Locator("text=Consult reports").WaitFor()
	assert.NoError(t, err)

	err = browser.Close()
	assert.NoError(t, err)

	err = pw.Stop()
	assert.NoError(t, err)

}
