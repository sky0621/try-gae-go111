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
	var (
		connectionName = os.Getenv("CLOUDSQL_CONNECTION_NAME")
		user           = os.Getenv("CLOUDSQL_USER")
		password       = os.Getenv("CLOUDSQL_PASSWORD")
		database       = os.Getenv("CLOUDSQL_DATABASE")
	)
	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@unix(/cloudsql/%s)/%s?parseTime=True", user, password, connectionName, database))
	if err != nil {
		panic(err)
	}
	defer db.Close()
	db.LogMode(true)

	e := echo.New()
	e.POST("/users", func(c echo.Context) error {

		e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", os.Getenv("PORT"))))
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
}

type User struct {
	ID   string `gorm:"column:id"`
	Name string `gorm:"column:name"`
	Mail string `gorm:"column:mail"`
}

func (u *User) TableName() string {
	return "user"
}
