package services

import (
	"bytes"
	"fmt"
	"os"
	"text/template"

	"github.com/asakshat/go-event-booking/internal/models"
	"gopkg.in/gomail.v2"
)

func SendGoMail(templatePath string, ticket models.TicketDetails) error {
	emailHost := os.Getenv("EMAIL")
	passwordHost := os.Getenv("PASSWORD")
	var body bytes.Buffer
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("EMAIL: ", emailHost)
	err = t.Execute(&body, ticket)
	if err != nil {
		fmt.Println(err)
		return err
	}
	m := gomail.NewMessage()
	m.SetHeader("From", emailHost)
	m.SetHeader("To", ticket.Email)
	m.SetHeader("Subject", "Ticket Purchase Confirmed!")
	m.SetBody("text/html", body.String())
	m.Embed(ticket.QRCodePath)

	d := gomail.NewDialer("smtp.gmail.com", 587, emailHost, passwordHost)

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
	return nil
}

func SendVerificationEmail(verifyemail models.EmailVerification) error {
	emailHost := os.Getenv("EMAIL")
	passwordHost := os.Getenv("PASSWORD")

	type EmailData struct {
		URL string
	}

	var body bytes.Buffer
	t, err := template.ParseFiles("./html/verify_email.html")
	if err != nil {
		fmt.Println(err)
		return err
	}
	emailData := EmailData{
		URL: fmt.Sprintf("http://localhost:8080/verify-email?token=%s", verifyemail.Token),
	}
	err = t.Execute(&body, emailData)
	if err != nil {
		fmt.Println(err)
		return err
	}
	m := gomail.NewMessage()
	m.SetHeader("From", emailHost)
	m.SetHeader("To", verifyemail.Email)
	m.SetHeader("Subject", "Email Verification")
	m.SetBody("text/html", body.String())

	d := gomail.NewDialer("smtp.gmail.com", 587, emailHost, passwordHost)
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
	return nil

}

func PasswordResetMail(templatePath string, email string, token string) {
	emailHost := os.Getenv("EMAIL")
	passwordHost := os.Getenv("PASSWORD")

	var body bytes.Buffer
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = t.Execute(&body, token)
	if err != nil {
		fmt.Println(err)
	}
	m := gomail.NewMessage()
	m.SetHeader("From", emailHost)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Password Reset")
	m.SetBody("text/html", body.String())

	d := gomail.NewDialer("smtp.gmail.com", 587, emailHost, passwordHost)

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}
