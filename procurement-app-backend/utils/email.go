package utils

import (
	"log"

	"gopkg.in/gomail.v2"
)

func SendEmail(to string, subject string, body string) {
	m := gomail.NewMessage()
	m.SetHeader("From", "your-email@example.com")
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer("smtp.example.com", 587, "your-email@example.com", "your-email-password")

	if err := d.DialAndSend(m); err != nil {
		log.Println(err)
	}
}
