package models

import (
	_ "github.com/lib/pq"
	"github.com/jinzhu/gorm"

	"fmt"
	"reflect"
	"errors"
)

type User struct {
	Name		string `gorm:"column:name;type:varchar(64);primary_key"`
	Password	string `gorm:"column:password;type:char(64)"`
	Session		string `gorm:"column:session;type:char(64)"`
}

var db *gorm.DB

func init() {
	var err error
	//db, err = gorm.Open("postgres", "golang:golang@/book_db?charset=utf8?sslmode=require")
	db, err = gorm.Open("postgres", "user=golang password=golang dbname=book_db sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer fmt.Println("close")
	defer db.Close()
	defer fmt.Println("now close")
	fmt.Println("succeed")
}

func open() {
	var err error
	//db, err = gorm.Open("postgres", "golang:golang@/book_db?charset=utf8?sslmode=require")
	db, err = gorm.Open("postgres", "user=golang password=golang dbname=book_db sslmode=disable")
	if err != nil {
		panic(err)
	}
}

// ユーザを作成するかエラーを返す
func CreateOrTop(user *User) error {
	open()
	defer db.Close()
	err := check(user)
	if err == nil {
		err = errors.New("impossible")
		return err
	}
	db.Create(user)
	return nil
}

// ユーザが存在するか確認する
func check(user *User) error {
	err := db.First(user, "name = ?", user.Name).Error
	return err
}

func Find(user *User) {
	open()
	defer db.Close()
	u := []User{}
	db.Find(&u, user)
	fmt.Println(u)

}

func Show() {
	fmt.Println(reflect.TypeOf(db))
}
