package main

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/labstack/echo"
)

func main() {
	e := echo.New()
	rand.Seed(time.Now().UnixNano())

	e.GET("/", func(c echo.Context) error {
		t := rand.Intn(1000)
		time.Sleep(time.Duration(t) * time.Millisecond)
		return c.JSON(http.StatusOK, "OK")
	})

	e.Logger.Fatal(e.Start(":8080"))
}
