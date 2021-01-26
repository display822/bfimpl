package services

import (
	"bytes"
	"html/template"
	"time"

	"github.com/astaxie/beego"
	"gopkg.in/gomail.v2"

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
func EmailExpenseApproved(mailTo string, id uint, name string, time time.Time) {
	subject := "报销审核通知"
	var body bytes.Buffer
	t, _ := template.ParseFiles("static/mail/approved.html")
	t.Execute(&body, struct {
		ID   uint
		Name string
		Time string
	}{
		ID:   id,
		Name: name,
		Time: time.Format("2006/01/02"),
	})
	sendMail(mailTo, subject, body.String())
}

// EmailExpenseRejectedUp 通过报销线上审批驳回通知
func EmailExpenseRejectedUp(mailTo string, name string, time time.Time, otp string) {
	subject := "报销审核通知"
	var body bytes.Buffer
	t, _ := template.ParseFiles("static/mail/rejectedUp.html")
	t.Execute(&body, struct {
		Name string
		Time string
		OTP  string
	}{
		Name: name,
		Time: time.Format("2006/01/02"),
		OTP:  otp,
	})
	sendMail(mailTo, subject, body.String())
}

// EmailExpenseRejectedDown 通过报销线下审批驳回通知
func EmailExpenseRejectedDown(mailTo string, name string, time time.Time, otp string) {
	subject := "报销支付通知"
	var body bytes.Buffer
	t, _ := template.ParseFiles("static/mail/rejectedDown.html")
	t.Execute(&body, struct {
		Name string
		Time string
		OTP  string
	}{
		Name: name,
		Time: time.Format("2006/01/02"),
		OTP:  otp,
	})
	sendMail(mailTo, subject, body.String())
}

// EmailExpensePaid 通过报销审批支付通知
func EmailExpensePaid(mailTo string, name string, expenseSummary float64, Acc string, time time.Time) {
	subject := "报销支付通知"
	var body bytes.Buffer
	t, _ := template.ParseFiles("static/mail/paid.html")
	t.Execute(&body, struct {
		Name           string
		ExpenseSummary float64
		Acc            string
		Time           string
	}{
		Name:           name,
		ExpenseSummary: expenseSummary,
		Acc:            Acc,
		Time:           time.Format("2006/01/02"),
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
