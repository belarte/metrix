package handlers

import (
	"github.com/belarte/metrix/views"
	"github.com/labstack/echo/v4"
)

func HomeHandler(c echo.Context) error {
	home := views.HomePage()
	return render(c, home)
}
