package main

import (
	"fmt"
	"os"

	"github.com/labstack/echo"
)

func main() {
	e := echo.New()

	// Google App Engineではデプロイ時にポートが決まり、環境変数「PORT」に接続ポートがセットされる。
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", os.Getenv("PORT"))))
}
