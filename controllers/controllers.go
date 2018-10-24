package controllers

import (
	"github.com/gin-gonic/gin"

	"net/http"
	"fmt"
	"reflect"
	"crypto/sha256"
	"io"

	m "book_manage/models"
)

var domain string = "localhost"

func GetLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "login.tmpl", gin.H{
		"title": "Login page",
	})

}

func PostLogin(c *gin.Context) {
	username := c.PostForm("username")
	pass := c.PostForm("password")
	c.JSON(http.StatusOK, gin.H{
		"username": username,
		"password": pass,
	})
}

func Startpage(c *gin.Context) {
	fmt.Println(c.ContentType())
	c.HTML(http.StatusOK, "top.tmpl", gin.H{
		"title": "Top page",
	})
	fmt.Println(c.ContentType())
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

	user := m.User{Name: name, Password: password}

	err := m.CreateOrTop(&user)

	if err != nil {
		fmt.Println(err)
		fmt.Println("impossible")
		c.Redirect(http.StatusMovedPermanently, "/")
	} else {
		c.SetCookie("user", name, 1800, "", domain, false, false)
		c.Redirect(http.StatusMovedPermanently, "/home")
	}
}

func GetHome(c *gin.Context) {
	user, err := c.Cookie("user")
	if err != nil {
		fmt.Println(err)
		user = "impossible"
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}

func SetCookie(c *gin.Context) {
	name := "username"
	value := "makura"
	//c.SetCookie(name, value, 1800, "", "localhost", false, false)
	c.SetCookie(name, value, 1800, "", "localhost", false, false)
	c.JSON(http.StatusOK, gin.H{"name": name, "value": value})
}

func Test(c *gin.Context) {
	c.Redirect(http.StatusMovedPermanently, "https://qiita.com")
	c.JSON(http.StatusOK, gin.H{"successed": "successed"})
}
