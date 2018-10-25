package controllers

import (
	"github.com/gin-gonic/gin"

	"net/http"
	"fmt"
	"crypto/sha256"
	"io"
	"errors"

	m "book_manage/models"
)

var domain string = "localhost"

func GetLogin(c *gin.Context) {
	var err error = errors.New("tmp")
	user, errU := c.Cookie("user")
	secret, errS := c.Cookie("secret")
	// sessionのチェック
	if errU == nil && errS == nil {
		err = m.CheckSessionLogin(user, secret)
	}
	if err == nil {
		c.Redirect(http.StatusMovedPermanently, "/home")
	} else {
		c.HTML(http.StatusOK, "login.tmpl", gin.H{
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
	err := m.CheckLogin(username, pass)

	if err != nil {
		c.HTML(http.StatusOK, "login.tmpl", gin.H{
			"title": "Login page",
		})
	} else {
		// sessionを作成
		io.WriteString(h, pass)
		session := fmt.Sprintf("%x", h.Sum(nil))
		user := m.User{Name: username, Session: session}

		// sessionの更新
		err = m.SetSession(&user)
		fmt.Println(err)
		c.SetCookie("user", username, 0, "/", domain, false, false)
		c.SetCookie("secret", session, 0, "/", domain, false, false)
		c.Redirect(http.StatusMovedPermanently, "/home")
	}
}

func Top(c *gin.Context) {
	c.HTML(http.StatusOK, "top.tmpl", gin.H{
		"title": "Top page",
	})
}

func GetSignUp(c *gin.Context) {
	c.HTML(http.StatusOK, "signup.tmpl", gin.H{
		"title": "SignUp",
	})
}

func PostSignUp(c *gin.Context) {
	name := c.PostForm("name")
	password := c.PostForm("password")

	//暗号化
	h := sha256.New()
	io.WriteString(h, password)
	password = fmt.Sprintf("%x", h.Sum(nil))

	io.WriteString(h, name)
	session := fmt.Sprintf("%x", h.Sum(nil))

	user := m.User{Name: name, Password: password, Session: session}

	// ユーザが存在するかチェック
	err := m.CreateOrTop(&user)

	if err == nil {
		c.SetCookie("user", name, 0, "", domain, false, false)
		c.SetCookie("secret", session, 0, "", domain, false, false)

		m.CreateTable(name+"booktable")
		c.Redirect(http.StatusMovedPermanently, "/home")
	} else {
		fmt.Println(err)
		c.HTML(http.StatusOK, "signup.tmpl", gin.H{
			"title": "SignUp",
		})
	}
}


func GetHome(c *gin.Context) {
	name, err := c.Cookie("user")
	if err != nil {
		fmt.Println(err)
		c.Redirect(http.StatusMovedPermanently, "/login")
		c.Abort()
	}
	session, err := c.Cookie("secret")
	if err != nil {
		fmt.Println(err)
		c.Redirect(http.StatusMovedPermanently, "/login")
		c.Abort()
	}
	err = m.CheckSessionLogin(name, session)
	if err != nil {
		fmt.Println(err)
		c.Redirect(http.StatusMovedPermanently, "/login")
		c.Abort()
	} else {
		book := m.GetBookData(name)
		c.HTML(http.StatusOK, "home.tmpl", gin.H{"Book": book, "title": "Home", "name": name})
	}
}

// ログアウト
func Logout(c *gin.Context) {
	c.SetCookie("user", "", -1, "/", domain, false, false)
	c.SetCookie("secret", "", -1, "/", domain, false, false)
	c.Redirect(http.StatusMovedPermanently, "/")
}

func SessionDelete() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.SetCookie("user", "", -1, "/", domain, false, false)
		c.SetCookie("secret", "", -1, "/", domain, false, false)
		c.Next()
	}
}



func SetCookie(c *gin.Context) {
	name := "username"
	value := "makura"
	c.SetCookie("user", value, 1800, "/", "localhost", false, false)
	c.SetCookie("secret", value, 1800, "/", "localhost", false, false)
	name, _ = c.Cookie("user")
	value, _ = c.Cookie("secret")
	c.JSON(http.StatusOK, gin.H{"name": name, "value": value})
}

func Test(c *gin.Context) {
	c.Redirect(http.StatusMovedPermanently, "https://qiita.com")
	c.JSON(http.StatusOK, gin.H{"successed": "successed"})
}
