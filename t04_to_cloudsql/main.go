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
	db, err := gorm.Open("mysql",
		fmt.Sprintf("root:%s@unix(/cloudsql/%s)/fs14db01?parseTime=True",
			os.Getenv("PASS"), os.Getenv("CONN")))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			panic(err)
		}
	}()
	// --------------------------------------------------------------
	// Pattern 1. Unlimited-open.
	//db.DB().SetMaxIdleConns(0)
	//db.DB().SetMaxOpenConns(0)

	// Pattern 2. < AppEngineInstance to CloudSQLInstance open.
	//db.DB().SetMaxIdleConns(0)
	//db.DB().SetMaxOpenConns(95)

	// Pattern 3. > AppEngineInstance to CloudSQLInstance open.
	//db.DB().SetMaxIdleConns(0)
	//db.DB().SetMaxOpenConns(200)
	// --------------------------------------------------------------

	e := echo.New()
	e.POST("/user", func(c echo.Context) error {
		id := uuid.New().String()
		u := &User{
			ID:   id,
			Name: fmt.Sprintf("user-%s", id),
			Mail: fmt.Sprintf("mail-%s@example.com", id),
		}
		if err := db.Save(u).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		}
		return c.JSON(http.StatusOK, "OK")
	})
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", os.Getenv("PORT"))))
}

type User struct {
	ID   string `gorm:"column:id"`
	Name string `gorm:"column:name"`
	Mail string `gorm:"column:mail"`
}

func (u *User) TableName() string {
	return "user"
}
