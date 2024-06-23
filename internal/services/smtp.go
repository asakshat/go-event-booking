package services

import (
	"bytes"
	"fmt"
	"os"
	"text/template"

	"github.com/asakshat/go-event-booking/internal/models"
	"gopkg.in/gomail.v2"
)

func SendGoMail(templatePath string, ticket models.TicketDetails) {
	emailHost := os.Getenv("EMAIL")
	passwordHost := os.Getenv("PASSWORD")

	var body bytes.Buffer
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = t.Execute(&body, ticket)
	if err != nil {
		fmt.Println(err)
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
}
