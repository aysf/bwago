package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aysf/bwago/internal/models"
	mail "github.com/xhit/go-simple-mail/v2"
)

func listenForMail() {
	go func() {
		for {
			m := <-app.MailChan
			sendMsg(m)
		}
	}()
}

func sendMsg(m models.MailData) {
	server := mail.NewSMTPClient()
	server.Host = "localhost"
	server.Port = 1025
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	client, err := server.Connect()
	if err != nil {
		errorlog.Println("error connecting mail server:", err)
	}
	email := mail.NewMSG()
	email.SetFrom(m.From).AddTo(m.To).SetSubject(m.Subject)
	if m.Template == "" {
		email.SetBody(mail.TextHTML, m.Content)
	} else {
		data, err := os.ReadFile(fmt.Sprintf("./email-templates/%s", m.Template))
		if err != nil {
			app.ErrorLog.Println("error read email template")
		}
		mailTemplate := string(data)
		msgToSend := strings.Replace(mailTemplate, "[%body%]", m.Content, 1)
		email.SetBody(mail.TextHTML, msgToSend)
	}

	err = email.Send(client)
	if err != nil {
		errorlog.Println("error sending email:", err)
	} else {
		infolog.Println("email sent")
	}

}
