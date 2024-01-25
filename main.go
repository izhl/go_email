package main

import (
	"fmt"
	"net/http"
	"net/smtp"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.POST("/to_email", send_email)
	r.Run(":8086") // 监听并在 0.0.0.0:8080 上启动服务
}

// 定义接收数据的结构体
type Email struct {
	// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
	To_email   []string `form:"email" json:"email" uri:"email" xml:"email" binding:"required"`
	Title      string   `form:"title" json:"title" uri:"title" xml:"title" binding:"required"`
	Content    string   `form:"content" json:"content" uri:"content" xml:"content" binding:"required"`
	Email_name string   `form:"email_name" json:"email_name" uri:"email_name" xml:"email_name" binding:"required"`
	Email_pass string   `form:"email_pass" json:"email_pass" uri:"email_pass" xml:"email_pass" binding:"required"`
	Email_host string   `form:"email_host" json:"email_host" uri:"email_host" xml:"email_host" binding:"required"`
	Email_port string   `form:"email_port" json:"email_port" uri:"email_port" xml:"email_port" binding:"required"`
	Email_user string   `form:"email_user" json:"email_user" uri:"email_user" xml:"email_user" binding:"required"`
}

func send_email(c *gin.Context) {
	// 创建在 goroutine 中使用的副本
	cCp := c.Copy()
	// 声明接收的变量
	var email Email
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := cCp.ShouldBindJSON(&email); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	to_email := email.To_email
	title := email.Title
	content_type := "Content-Type: text/html; charset=UTF-8"
	// content_arr := email.Content.Split(`src=\"`)
	email_content := `
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="iso-8859-15">
			<title>MMOGA POWER</title>
		</head>
		<body>
			` + fmt.Sprintf("%s", email.Content) +
		`</body>
		</html>`

	for _, t := range to_email {
		go send_to_email(t, title, content_type, email_content, email.Email_name, email.Email_pass, email.Email_host, email.Email_port, email.Email_user)
	}

	c.JSON(200, gin.H{
		"message": to_email,
	})
}

func send_to_email(to_email, title, content_type, content, email_name, email_pass, email_host, email_port, email_user string) error {

	auth := smtp.PlainAuth("", email_name, email_pass, email_port)
	msg := []byte("To: " + to_email + "\r\nFrom: " + email_user + "<" + email_name + ">" + "\r\nsubject: " + title + "\r\n" + content_type + "\r\n\r\n" + content)
	send_to := strings.Split(to_email, ";")
	err := smtp.SendMail(email_host+":"+email_port, auth, email_name, send_to, msg)
	fmt.Println(err)
	return err
}
