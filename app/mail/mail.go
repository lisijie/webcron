package mail

import (
	"github.com/astaxie/beego"
	"github.com/lisijie/gomail"
	"time"
)

var (
	sendCh    chan *gomail.Message
	smtpHost  string
	smtpPort  int
	mailFrom  string
	username  string
	password  string
	queueSize int
)

func init() {
	queueSize, _ = beego.AppConfig.Int("mail.queue_size")
	smtpHost = beego.AppConfig.String("mail.smtp")
	smtpPort, _ = beego.AppConfig.Int("mail.port")
	username = beego.AppConfig.String("mail.user")
	password = beego.AppConfig.String("mail.password")
	mailFrom = beego.AppConfig.String("mail.from")

	if smtpPort == 0 {
		smtpPort = 25
	}

	sendCh = make(chan *gomail.Message, queueSize)

	go func() {
		d := gomail.NewPlainDialer(smtpHost, smtpPort, username, password)
		var s gomail.SendCloser
		var err error
		open := false
		for {
			select {
			case m, ok := <-sendCh:
				if !ok {
					return
				}
				if !open {
					if s, err = d.Dial(); err != nil {
						beego.Error(err)
						continue
					}
					open = true
				}
				if err := gomail.Send(s, m); err != nil {
					beego.Error("SendMail:", err.Error())
				}
			case <-time.After(30 * time.Second):
				if open {
					if err := s.Close(); err != nil {
						beego.Error(err)
					}
					open = false
				}
			}
		}
	}()
}

func SendMail(address, name, subject, content string) bool {
	msg := gomail.NewMessage()
	msg.SetHeader("From", mailFrom)
	msg.SetAddressHeader("To", address, name)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", content)

	select {
	case sendCh <- msg:
		return true
	case <-time.After(time.Second * 3):
		return false
	}
}
