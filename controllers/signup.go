package controllers

import (
	"github.com/gin-gonic/gin"

	"fmt"
	"net/http"
	"crypto/sha256"
	"io"
	"os"
	"errors"

	m "book_manage/models"
)

func GetSignUp(c *gin.Context) {
	css := []string{"assets/css/signup/signup.css"}
	c.HTML(http.StatusOK, "signup.tmpl", gin.H{
		"title": "SignUp",
		"CSS": css,
	})
}

func PostSignUp(c *gin.Context) {
	name := c.PostForm("name")
	var err error
	if name == "" {
		c.Redirect(http.StatusSeeOther, "/signup")
		err = errors.New("login ID lose")
		return
	}

	password := c.PostForm("password")

	//暗号化
	h := sha256.New()
	io.WriteString(h, password)
	password = fmt.Sprintf("%x", h.Sum(nil))

	// session作成
	io.WriteString(h, name)
	session := fmt.Sprintf("%x", h.Sum(nil))

	user := m.User{Name: name, Password: password, Session: session}

	if err == nil {
		// ユーザが存在するかチェック
		err = m.CreateOrTop(&user)
	}

	if err == nil {
		c.SetCookie("user", name, 0, "", domain, false, false)
		c.SetCookie("secret", session, 0, "", domain, false, false)

		m.CreateTable(name+"booktable")
		// ユーザ毎に画像を保存するパスを作成
		path := "./assets/img/"+name
		err = os.MkdirAll(path, 0750)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
		c.Redirect(http.StatusSeeOther, "/home")
	} else {
		c.Redirect(http.StatusSeeOther, "/signup")
		fmt.Println(err)
	}
}
