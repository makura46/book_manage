package models

import (
	_ "github.com/lib/pq"
	"github.com/jinzhu/gorm"

	"fmt"
	"reflect"
	"errors"
)

type User struct {
	Name		string `gorm:"column:name;type:varchar(64);primary_key;NOT NULL"`
	Password	string `gorm:"column:password;type:char(64)"`
	Session		string `gorm:"column:session;type:char(64)"`
}

type BookTable struct {
	Id			int		`gorm:"column:id;primary_key;NOT NULL;AUTO_INCREMENT"`
	ImgPath		string	`gorm:"column:imgpath;type:varchar(256)"`
	Name		string	`gorm:"column:bookname;type:varchar(256)"`
	Read		int		`gorm:"column:read;type:int"`
	Have		int		`gorm:"column:have;type:int"`
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

// データベースに接続する
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
	// すでに存在する
	if err == nil {
		err = errors.New("impossible")
		return err
	}

	// ユーザーを作成する
	db.Create(user)
	return nil
}

// ユーザがログイン状態かuserNameとsessionから確認する
func CheckSessionLogin(userName, session string) error {
	open()
	defer db.Close()
	err := db.First(&User{}, "name = ? AND session = ?", userName, session).Error
	return err
}

// ログインチェック
// ユーザ登録がされているか確認する
func CheckLogin(name, pass string) error {
	open()
	defer db.Close()

	err := db.First(&User{}, "name = ? AND password = ?", name, pass).Error
	fmt.Println(err)
	return err
}

// テーブルを作成する
func CreateTable(bookTableName string) {
	open()
	defer db.Close()
	db.Table(bookTableName).CreateTable(&BookTable{})
}

// テーブルからユーザの本情報を取得する
func AllGetBookData(name string) []BookTable {
	open()
	defer db.Close()
	tableName := name+"booktable"
	var book []BookTable
	db.Table(tableName).Order("id asc").Find(&book)
	return book
}

//sessionを更新する
func SetSession(u *User) error {
	open()
	defer db.Close()
	err := db.Model(u).Updates(u).Error
	return err
}

// 本の情報をidから取得する
func GetBookData(tableName string, id int) (*BookTable, error) {
	open()
	defer db.Close()
	bookData := BookTable{Id: id}
	err := db.Table(tableName).First(&bookData).Error
	fmt.Println(bookData)
	return &bookData, err
}

// 本の情報を変更する
func UpdateRecord(tableName string, book *BookTable, id int) error {
	open()
	defer db.Close()
	fmt.Println(*book)
	err := db.Table(tableName).Where("id = ?", id).Updates(book).Error
	return err
}

func AddRecord(tableName string, book *BookTable) error {
	open()
	defer db.Close()
	err := db.Table(tableName).Create(book).Error
	return err
}

func DeleteBook(tableName string, book *BookTable) error {
	open()
	defer db.Close()
	err := db.Table(tableName).Delete(book).Error
	return err
}


func Find(user *User) {
	open()
	defer db.Close()
	u := []User{}
	db.Find(&u, user)
	fmt.Println(u)

}

// ユーザが存在するか確認する
func check(user *User) error {
	err := db.First(&User{}, "name = ?", user.Name).Error
	return err
}

