package main

import (
	"net/http"
	"time"

	"github.com/labstack/echo"
)

func main() {
	e := echo.New()
	g := e.Group("/api/v1")

	// デフォルトスケーリング設定用（ドキュメント上、１HTTPリクエストは「60秒」以内に終える必要があると記載あり。）
	g.POST("/automatic", func(c echo.Context) error {
		time.Sleep(10 * time.Minute) // 10分スリープして処理が返ってくるか？
		return c.JSON(http.StatusOK, "OK")
	})

	// ベーシックスケーリング設定用（こちらは最大 24時間までOK）
	g.GET("/basic", func(c echo.Context) error {
		time.Sleep(10 * time.Minute) // 10分スリープして処理が返ってくるか？
		return c.JSON(http.StatusOK, "OK")
	})

	e.Logger.Fatal(e.Start(":8080"))
}
