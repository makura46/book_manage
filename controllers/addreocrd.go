package controllers

import (
	"github.com/gin-gonic/gin"

	"net/http"
	"fmt"
	"os"
	"io"
	"strconv"

	m "book_manage/models"
)

func AddRecord(c *gin.Context) {
	user, secret, err := sessionCheck(c, "user", "secret")
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/login")
	}
	err = m.CheckSessionLogin(user, secret)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/login")
	}
	c.HTML(http.StatusOK, "addrecord.tmpl", gin.H{"title": "AddRecord", "user": user})
}

func PostRecord(c *gin.Context) {
	user, secret, err := sessionCheck(c, "user", "secret")
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/login")
	}
	err = m.CheckSessionLogin(user, secret)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/login")
	}
	var imgPath string
	imgTmp, err := c.FormFile("img")
	if err != nil {
		imgPath = "./assets/img/No_img.jpg"
	} else {
		// ファイル名の重複を考えない
		path := "./assets/img/"+user+"/"+imgTmp.Filename
		/*
		for {
			err := isFileExist(path)
			if err == nil {
				break
			}
		}
		*/
		imgPath = path
		var file *os.File
		file, err = os.Create(imgPath)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
		defer file.Close()
		img, e := imgTmp.Open()
		if e != nil {
			panic(err)
		}
		defer img.Close()
		_, err = io.Copy(file, img)
		if err != nil {
			panic(err)
		}
	}
	title := c.PostForm("title")
	readString := c.PostForm("read")
	haveString := c.PostForm("have")
	stopFlag := false
	if title == "" {
		c.HTML(http.StatusOK, "addrecord.tmpl", gin.H{"title": "AddRecord", "user": user})
		stopFlag = true
	}
	var read, have int
	read, err = strconv.Atoi(readString)
	if err != nil && !stopFlag {
		c.HTML(http.StatusOK, "addrecord.tmpl", gin.H{"title": "AddRecord", "user": user})
		stopFlag = true
	}
	have, err = strconv.Atoi(haveString)
	if err != nil && !stopFlag {
		c.HTML(http.StatusOK, "addrecord.tmpl", gin.H{"title": "AddRecord", "user": user})
		stopFlag = true
	}
	if !stopFlag {
		book := m.BookTable{ImgPath: imgPath, Name: title, Read: read, Have: have}
		m.AddRecord(user+"booktable", &book)
		c.Redirect(http.StatusSeeOther, "/home")
	}
}
