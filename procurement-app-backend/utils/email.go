package utils

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
)

func init() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func SendEmail(to, subject, body string) {
	// Get environment variables
	email := os.Getenv("OUTLOOK_EMAIL")
	password := os.Getenv("OUTLOOK_PASSWORD")

	m := gomail.NewMessage()
	m.SetHeader("From", email)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer("smtp-mail.outlook.com", 587, email, password)

	if err := d.DialAndSend(m); err != nil {
		log.Println("Error sending email:", err)
	}
}

func OrderCreatedEmail(to, orderID string, totalAmount float64) {
	subject := "Order Created"
	body := fmt.Sprintf(`
		<h1>Order Created</h1>
		<p>Dear %s,</p>
		<p>Your order with ID %s has been successfully created.</p>
		<p>Total Amount: $%.2f</p>
		<p>Thank you for your purchase!</p>
	`, to, orderID, totalAmount)

	SendEmail(to, subject, body)
}
