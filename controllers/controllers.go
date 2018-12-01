package controllers

import (
	"github.com/gin-gonic/gin"

	"net/http"
	"fmt"
	"strconv"

	m "book_manage/models"
)

func Top(c *gin.Context) {
	c.HTML(http.StatusOK, "top.tmpl", gin.H{
		"title": "Top page",
	})
}

// ログアウト
func Logout(c *gin.Context) {
	c.SetCookie("user", "", -1, "/", domain, false, false)
	c.SetCookie("secret", "", -1, "/", domain, false, false)
	c.Redirect(http.StatusSeeOther, "/")
}


func Delete(c *gin.Context) {
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
	idString := c.PostForm("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/home")
		return
	}
	book := m.BookTable{Id: id}
	err = m.DeleteBook(user+"booktable", &book)
	if err != nil {
		fmt.Println(err)
	}
	c.Redirect(http.StatusSeeOther, "/home")
}


func SessionDelete() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.SetCookie("user", "", -1, "/", domain, false, false)
		c.SetCookie("secret", "", -1, "/", domain, false, false)
		c.Next()
	}
}

func sessionCheck(c *gin.Context, nameU, nameS string) (valueU, valueS string, err error) {
	valueU, err = c.Cookie(nameU)
	if err != nil {
		return "", "", err
	}
	valueS, err = c.Cookie(nameS)
	if err != nil {
		return "", "", err
	}
	return valueU, valueS, err
}

