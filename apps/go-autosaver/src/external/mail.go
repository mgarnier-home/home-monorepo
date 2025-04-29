package external

import (
	"crypto/tls"

	"gopkg.in/gomail.v2"

	"mgarnier11.fr/go/go-autosaver/config"
)

func SendMail(
	mailConfig *config.MailConfig,
	to,
	subject,
	body string,
) error {

	dialer := gomail.NewDialer(mailConfig.Host, mailConfig.Port, mailConfig.Login, mailConfig.Password)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	message := gomail.NewMessage()
	message.SetHeader("From", mailConfig.Login)
	message.SetHeader("To", to)
	message.SetHeader("Subject", subject)
	message.SetBody("text/plain", body)

	return dialer.DialAndSend(message)

}
