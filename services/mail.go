package services

import (
	"bytes"
	"github.com/astaxie/beego"
	"gopkg.in/gomail.v2"
	"html/template"

	"bfimpl/services/log"
)

var (
	dialer   *gomail.Dialer
	mailHost string
	mailPort int
	mailUser string
	mailPass string
)

func init() {
	mailHost = beego.AppConfig.String("MailHost")
	mailPort, _ = beego.AppConfig.Int("MailPort")
	mailUser = beego.AppConfig.String("MailUser")
	mailPass = beego.AppConfig.String("MailPass")
}

func MailInit() {
	dialer = gomail.NewDialer(mailHost, mailPort, mailUser, mailPass)
}

// EmailExpenseApproved 通过报销审批通过通知
func EmailExpenseApproved(mailTo string) {
	subject := "Approved"
	body := "Approved"
	sendMail(mailTo, subject, body)
}

// EmailExpenseRejectedUp 通过报销线上审批驳回通知
func EmailExpenseRejectedUp(mailTo string) {
	subject := "RejectedUp"
	body := "RejectedUp"
	sendMail(mailTo, subject, body)
}

// EmailExpenseRejectedDown 通过报销线下审批驳回通知
func EmailExpenseRejectedDown(mailTo string) {
	subject := "RejectedDown"
	body := "RejectedDown"
	sendMail(mailTo, subject, body)
}

// EmailExpensePaid 通过报销审批支付通知
func EmailExpensePaid(mailTo string) {
	subject := "paid"
	var body bytes.Buffer
	t, _ := template.ParseFiles("static/mail/paid.html")
	t.Execute(&body, struct {
		Name           string
		ExpenseSummary float64
		Acc            string
		Time           string
	}{
		Name:           "yi.zhang",
		ExpenseSummary: 123111.12,
		Acc:            "231df12321312312312321",
		Time:           "2020/05/06",
	})
	sendMail(mailTo, subject, body.String())
}

func sendMail(mailTo string, subject string, body string) {
	m := gomail.NewMessage()
	//m.SetHeader("From", m.FormatAddress(mailUser, "broadfun")) //这种方式可以添加别名，即“XX官方”
	m.SetHeader("From", mailUser)
	m.SetHeader("To", mailTo)       //发送给用户
	m.SetHeader("Subject", subject) //设置邮件主题
	m.SetBody("text/html", body)    //设置邮件正文*/
	err := dialer.DialAndSend(m)
	if err != nil {
		log.GLogger.Error("sendMail err :%s", err.Error())
	}
}
