package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"

	"github.com/jinzhu/gorm"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
)

func main() {
	// --------------------------------------------------------------
	// GCP - CloudSQL への接続情報は環境変数から取得
	// --------------------------------------------------------------
	var (
		connectionName = os.Getenv("CLOUDSQL_CONNECTION_NAME")
		user           = os.Getenv("CLOUDSQL_USER")
		password       = os.Getenv("CLOUDSQL_PASSWORD")
		database       = os.Getenv("CLOUDSQL_DATABASE")
	)
	dataSource := fmt.Sprintf("%s:%s@unix(/cloudsql/%s)/%s?parseTime=True", user, password, connectionName, database)

	// ローカル環境（docker-composeでMySQLを起動）での動作確認用
	isLocal := os.Getenv("IS_LOCAL") != ""
	if isLocal {
		dataSource = "testuser:testpass@tcp(localhost:3306)/testdb?charset=utf8&parseTime=True&loc=Local"
	}

	// --------------------------------------------------------------
	// GCP - CloudSQL へ接続
	// --------------------------------------------------------------
	db, err := gorm.Open("mysql", dataSource)
	if err != nil {
		panic(err)
	}
	defer func() {
		if db != nil {
			err := db.Close()
			if err != nil {
				panic(err)
			}
		}
	}()

	db.LogMode(true)
	if err := db.DB().Ping(); err != nil {
		panic(err)
	}

	// チューニングポイント
	db.DB().SetMaxIdleConns(100)
	db.DB().SetMaxOpenConns(100)

	rand.Seed(time.Now().UnixNano())

	e := echo.New()

	// --------------------------------------------------------------
	// 各種WebAPIの定義
	// --------------------------------------------------------------
	// ユーザー登録
	e.POST("/users", func(c echo.Context) error {
		id := uuid.New().String()
		u := &User{
			ID:   id,
			Name: fmt.Sprintf("ユーザー%s", id),
			Mail: fmt.Sprintf("mail-%s@example.com", id),
		}
		if err := db.Save(u).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		}
		return c.JSON(http.StatusOK, "OK")
	})

	// 全ユーザー取得
	e.GET("/users", func(c echo.Context) error {
		var u *User
		if err := db.Find(&u).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		}
		return c.JSON(http.StatusOK, u)
	})

	port := fmt.Sprintf(":%s", os.Getenv("PORT"))
	fmt.Printf("[GWFGS] PORT: %s\n", port)
	if isLocal {
		port = ":8080"
	}
	e.Logger.Fatal(e.Start(port))
}

type User struct {
	ID   string `gorm:"column:id"`
	Name string `gorm:"column:name"`
	Mail string `gorm:"column:mail"`
}

func (u *User) TableName() string {
	return "user"
}
