package controllers

import (
	"github.com/gin-gonic/gin"

	"net/http"
	"fmt"
	"crypto/sha256"
	"io"
	"strconv"

	m "book_manage/models"
)

var domain string = "localhost"

func GetLogin(c *gin.Context) {
	user, secret, err := sessionCheck(c, "user", "secret")
	// sessionのチェック
	if err == nil {
		err = m.CheckSessionLogin(user, secret)
	}
	if err == nil {
		c.Redirect(http.StatusMovedPermanently, "/home")
		c.Abort()
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
	if name == "" {
		c.HTML(http.StatusOK, "signup.tmpl", gin.H{
			"title": "SignUp",
		})
	}
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
	}
	session, err := c.Cookie("secret")
	if err != nil {
		fmt.Println(err)
		c.Redirect(http.StatusMovedPermanently, "/login")
	}
	err = m.CheckSessionLogin(name, session)
	if err != nil {
		fmt.Println(err)
		c.Redirect(http.StatusMovedPermanently, "/login")
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

func AddRecord(c *gin.Context) {
	user, secret, err := sessionCheck(c, "user", "secret")
	if err != nil {
		c.Redirect(http.StatusMovedPermanently, "/login")
	}
	err = m.CheckSessionLogin(user, secret)
	if err != nil {
		c.Redirect(http.StatusMovedPermanently, "/login")
	}
	c.HTML(http.StatusOK, "addrecord.tmpl", gin.H{"title": "AddRecord", "user": user})
}

func PostRecord(c *gin.Context) {
	user, secret, err := sessionCheck(c, "user", "secret")
	if err != nil {
		c.Redirect(http.StatusMovedPermanently, "/login")
	}
	err = m.CheckSessionLogin(user, secret)
	if err != nil {
		c.Redirect(http.StatusMovedPermanently, "/login")
	}
	title := c.PostForm("title")
	readString := c.PostForm("read")
	haveString := c.PostForm("have")
	if title == "" {
		c.HTML(http.StatusOK, "addrecord.tmpl", gin.H{"title": "AddRecord", "user": user})
	}
	var read, have int
	read, err = strconv.Atoi(readString)
	if err != nil {
		c.HTML(http.StatusOK, "addrecord.tmpl", gin.H{"title": "AddRecord", "user": user})
	}
	have, err = strconv.Atoi(haveString)
	if err != nil {
		c.HTML(http.StatusOK, "addrecord.tmpl", gin.H{"title": "AddRecord", "user": user})
	}
	book := m.BookTable{Name: title, Read: read, Have: have}
	m.AddRecord(user+"booktable", &book)
	c.Redirect(http.StatusMovedPermanently, "/home")
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
