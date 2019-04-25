package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo"
)

func main() {
	e := echo.New()
	g := e.Group("/api/v1")

	// デフォルトスケーリング設定用（ドキュメント上、１HTTPリクエストは「60秒」以内に終える必要があると記載あり。）
	g.GET("/automatic/min10", func(c echo.Context) error {
		time.Sleep(10 * time.Minute)
		return c.JSON(http.StatusOK, "OK")
	})

	// ベーシックスケーリング設定用（こちらは最大 24時間までOK）
	g.GET("/basic/min10", func(c echo.Context) error {
		time.Sleep(10 * time.Minute)
		return c.JSON(http.StatusOK, "OK")
	})

	// ベーシックスケーリング設定のディスパッチ確認用
	// （即座に結果を返す処理をもって、想定したサービスにリクエストが振られていることを確認する）
	g.GET("/basic/now", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "OK")
	})

	// Google App Engineではデプロイ時にポートが決まり、環境変数「PORT」に接続ポートがセットされる。
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", os.Getenv("PORT"))))
}
