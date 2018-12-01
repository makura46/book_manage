package controllers

import (
	"fmt"
	"net/http"
	"io"
	"os"
	"strconv"
	"mime/multipart"

	m "book_manage/models"

	"github.com/gin-gonic/gin"
)

// 変更する本情報の入力を受け付ける
func PostEdit(c *gin.Context) {
	user, secret, err := sessionCheck(c, "user", "secret")
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}
	err = m.CheckSessionLogin(user, secret)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}
	idStr := c.PostForm("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/home")
		return
	}
	var bookData *m.BookTable
	bookData, err = m.GetBookData(user+"booktable", id)
	if err != nil {
		c.Redirect(http.StatusSeeOther,  "/home")
		return
	}
	fmt.Println("----show bookData----")
	fmt.Println(bookData)
	fmt.Println("---------------------")
	c.HTML(http.StatusOK, "edit.tmpl", gin.H{"title": "Change Page", "Book": bookData})
}

// 本の情報を変更する
func PostChange(c *gin.Context) {
	user, secret, err := sessionCheck(c, "user", "secret")
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}
	err = m.CheckSessionLogin(user, secret)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}
	idStr := c.PostForm("id")
	id, _ := strconv.Atoi(idStr)

	// 画像の実態を取得
	imgTmp, e := c.FormFile("img")
	var imgPath string
	// 画像があれば
	if e == nil {
		// 画像のパスを生成
		imgPath = "./assets/img/"+user+"/"+imgTmp.Filename
		// 画像をオープン
		var img multipart.File
		img, err = imgTmp.Open()
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
		defer img.Close()

		// サーバーに保存するファイルを作成
		var file *os.File
		file, err = os.Create(imgPath)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
		// サーバーに画像を保存する
		_, err = io.Copy(file, img)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
	} else {
		imgPath := c.PostForm("img")
		if imgPath == "" {
			imgPath = "./assets/No.jpg";
		}
	}
	name := c.PostForm("name")
	if name == "" {
		c.Redirect(http.StatusSeeOther, "/home")
		return
	}
	readStr := c.PostForm("read")
	read, err := strconv.Atoi(readStr)
	if err != nil {
		fmt.Println("error")
		c.Redirect(http.StatusSeeOther, "/home")
		return
	}
	haveStr := c.PostForm("have")
	have, err := strconv.Atoi(haveStr)
	if err != nil {
		fmt.Println("error have")
		c.Redirect(http.StatusSeeOther, "/home")
		return
	}
	//book := m.BookTable{Id: id, Name: name, Read: read, Have: have}
	book := m.BookTable{ImgPath: imgPath, Name: name, Read: read, Have: have}
	fmt.Println(book)
	err = m.UpdateRecord(user+"booktable", &book, id)
	if err != nil {
		fmt.Println(err)
	}
	c.Redirect(http.StatusSeeOther, "/home")
}
