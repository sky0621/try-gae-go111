package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/jinzhu/gorm"

	"github.com/google/uuid"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
)

func main() {
	// ローカル環境での実行時は、環境変数にセットしておく。例：「export IS_LOCAL='yes'」
	isLocal := os.Getenv("IS_LOCAL") != ""

	// Cloud SQL (ローカルでは MySQL) への接続
	db, closeDBFunc, err := connectDB(isLocal)
	if err != nil {
		os.Exit(-1)
	}
	defer closeDBFunc(db)

	// WebAPIサーバ起動
	e := echo.New()
	routing(e, db)

	port := fmt.Sprintf(":%s", os.Getenv("PORT"))
	fmt.Printf("[GWFGS] PORT: %s\n", port)
	if isLocal {
		port = ":8080"
	}
	e.Logger.Fatal(e.Start(port))
}

func createDatasource(isLocal bool) string {
	if isLocal {
		return "testuser:testpass@tcp(localhost:3306)/testdb?charset=utf8&parseTime=True&loc=Local"
	}

	// GCP - CloudSQL への接続情報は環境変数から取得
	var (
		connectionName = os.Getenv("CLOUDSQL_CONNECTION_NAME")
		user           = os.Getenv("CLOUDSQL_USER")
		password       = os.Getenv("CLOUDSQL_PASSWORD")
		database       = os.Getenv("CLOUDSQL_DATABASE")
	)

	return fmt.Sprintf("%s:%s@unix(/cloudsql/%s)/%s?parseTime=True", user, password, connectionName, database)
}

type closeDB func(db *gorm.DB)

func closeDBFunc(db *gorm.DB) {
	if db == nil {
		return
	}
	_ = db.Close()
}

func connectDB(isLocal bool) (*gorm.DB, closeDB, error) {
	db, err := gorm.Open("mysql", createDatasource(isLocal))
	if err != nil {
		return nil, nil, err
	}

	db.LogMode(true)
	if err := db.DB().Ping(); err != nil {
		return nil, nil, err
	}

	db.DB().SetMaxIdleConns(10000)

	return db, closeDBFunc, nil
}

func routing(e *echo.Echo, db *gorm.DB) {
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
}

type User struct {
	ID   string `gorm:"column:id"`
	Name string `gorm:"column:name"`
	Mail string `gorm:"column:mail"`
}

func (u *User) TableName() string {
	return "user"
}
