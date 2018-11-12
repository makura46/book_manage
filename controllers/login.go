package controllers

import (
	"github.com/gin-gonic/gin"

	"net/http"
	"fmt"
	"crypto/sha256"
	"io"

	m "book_manage/models"
)

func GetLogin(c *gin.Context) {
	user, secret, err := sessionCheck(c, "user", "secret")
	// sessionのチェック
	if err == nil {
		err = m.CheckSessionLogin(user,secret)
	}
	if err == nil {
		c.Redirect(http.StatusSeeOther,  "/home")
	} else {
		c.HTML(http.StatusOK, "login.tmpl", gin.H {
			"title": "Login page",
		})
	}
}

func PostLogin(c *gin.Context) {
	username := c.PostForm("username")
	pass := c.PostForm("password")

	h := sha256.New()
	io.WriteString(h, pass)
	pass = fmt.Sprintf("%x", h.Sum(nil))

	// ログインチェック
	// ユーザが登録されているか確認する
	// ユーザが登録されていればパスワードがあっているかも確認する
	err := m.CheckLogin(username, pass)

	if err != nil {
		c.HTML(http.StatusOK, "login.tmpl", gin.H{
			"title": "Login page",
		})
	} else {
		// sessionを作成
		io.WriteString(h, username)
		session := fmt.Sprintf("%x", h.Sum(nil))
		user := m.User{Name: username, Session: session}

		// sessionの更新
		err = m.SetSession(&user)
		fmt.Println(err)
		c.SetCookie("user", username, 0, "/", domain, false, false)
		c.SetCookie("secret", session, 0, "/", domain, false, false)
		c.Redirect(http.StatusSeeOther, "/home")
	}
}
