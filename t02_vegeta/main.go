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
		return c.JSON(http.StatusOK, "OK")
	})

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", os.Getenv("PORT"))))
}
