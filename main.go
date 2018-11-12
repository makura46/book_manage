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
	r.Static("/vendor", "./vendor")

	r.GET("/", c.Top)
	r.GET("/login", c.GetLogin)
	r.POST("/login", c.PostLogin)
	r.GET("/signup", c.GetSignUp)
	r.POST("/signup", c.PostSignUp)
	r.GET("/home", c.GetHome)
	r.GET("/logout", c.Logout)
	r.POST("/logout", c.Logout)
	r.GET("/addrecord", c.AddRecord)
	r.POST("/addrecord", c.PostRecord)
	r.POST("/delete", c.Delete)
	r.POST("/edit", c.PostEdit)
	r.POST("/change", c.PostChange)


	r.Run(":8080")
}
