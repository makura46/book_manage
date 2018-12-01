package controllers

import (
	"github.com/gin-gonic/gin"

	"net/http"
	"fmt"

	m "book_manage/models"
)

func GetHome(c *gin.Context) {
	name, err := c.Cookie("user")
	if err != nil {
		fmt.Println(err)
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}
	session, err := c.Cookie("secret")
	if err != nil {
		fmt.Println(err)
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}
	err = m.CheckSessionLogin(name, session)
	if err != nil {
		fmt.Println(err)
		c.Redirect(http.StatusSeeOther, "/login")
	} else {
		book := m.AllGetBookData(name)
		c.HTML(http.StatusOK, "home.tmpl", gin.H{"Book": book, "title": "Home", "name": name})
	}
}
