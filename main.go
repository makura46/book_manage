package main

import (
	"github.com/gin-gonic/gin"
	"net/http"

	c "book_manage/controllers"
)

func startpage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.tmpl", gin.H{
		"title": "Login page",
	})
}

func main() {
	r := gin.Default()

	r.LoadHTMLGlob("views/**/*")
	r.Static("/assets", "./assets")

	r.GET("/", c.Startpage)
	r.GET("/login", c.GetLogin)
	r.POST("/login", c.PostLogin)
	r.GET("/signup", c.GetSignUp)
	r.POST("/signup", c.PostSignUp)
	r.GET("/home", c.GetHome)
	r.GET("/test", c.Test)
	/*
	r.GET("/setCookie", c.SetCookie)
	r.GET("/getCookie", c.GetCookie)
	r.GET("deleteCookie", c.DeleteCookie)
	*/

	r.Run(":8080")
}
