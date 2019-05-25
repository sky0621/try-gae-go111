package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

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
	// Pattern 1
	db.DB().SetMaxIdleConns(0)
	db.DB().SetMaxOpenConns(0)

	// Pattern 2
	//db.DB().SetMaxIdleConns(0)
	//db.DB().SetMaxOpenConns(95)

	// Pattern 3
	//db.DB().SetMaxIdleConns(0)
	//db.DB().SetMaxOpenConns(200)

	// Pattern 4
	//db.DB().SetMaxIdleConns(95)
	//db.DB().SetMaxOpenConns(95)

	// Pattern 5
	//db.DB().SetMaxIdleConns(200)
	//db.DB().SetMaxOpenConns(200)

	// Pattern 5
	//db.DB().SetMaxIdleConns(1000)
	//db.DB().SetMaxOpenConns(1000)
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
		fmt.Println(time.Now())
		db.Raw("SELECT sleep(3)")
		fmt.Println(time.Now())

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
