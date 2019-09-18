package email

import (
	"fmt"
	"net/smtp"
	"strings"
)

type Email struct {
	From     string
	AuthCode string
	Host     string
	Port     int
	Sender   string
}

type EmailConf struct {
	From     string `json:"From"`
	AuthCode string `json:"AuthCode"`
	Host     string `json:"Host"`
	Port     int    `json:"Port"`
	Sender   string `json:"Sender"`
}

func NewEmail(conf EmailConf) *Email {
	return &Email{
		From:     conf.From,
		AuthCode: conf.AuthCode,
		Host:     conf.Host,
		Port:     conf.Port,
		Sender:   conf.Sender,
	}
}

// 发送
func (e *Email) Send(to, subject, content string) error {
	from := e.From

	auth := smtp.PlainAuth("", from, e.AuthCode, e.Host)

	headers := make(map[string]string)
	headers["From"] = e.Sender + "<" + from + ">"
	headers["To"] = to
	headers["Subject"] = subject
	headers["Content-Type"] = "text/html; charset=UTF-8"

	fmtCon := ""
	for key, value := range headers {
		fmtCon += fmt.Sprintf("%s: %s\r\n", key, value)
	}
	fmtCon += "\r\n" + content

	// contentType := "Content-Type: text/html; charset=UTF-8"

	sendTo := strings.Split(to, ";")
	// fmtCon := "To: " + to + "\r\nFrom: " + from + "\r\nSubject: " + subject + "\r\n" + contentType + "\r\n\r\n" + content

	addr := fmt.Sprintf("%s:%d", e.Host, e.Port)
	err := smtp.SendMail(addr, auth, from, sendTo, []byte(fmtCon))
	return err
}

// func SendEmail(to, subject, content string) {
// 	if err := email.Send(to, subject, content); err != nil {
// 		logger.Debugln(err)
// 		logger.Errorln("注册码发生失败")
// 	}
// }
