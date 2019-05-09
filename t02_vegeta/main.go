package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo"
)

func main() {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		if c.Request().Header.Get("API-KEY") != "3a193b40-9691-4f4f-841a-4417ff15e110" {
			return c.JSON(http.StatusForbidden, nil)
		}
		return c.JSON(http.StatusOK, "OK")
	})

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", os.Getenv("PORT"))))
}
