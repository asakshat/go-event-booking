package services

import (
	"bytes"
	"fmt"
	"os"
	"text/template"

	"gopkg.in/gomail.v2"
)

type TicketDetails struct {
	CustomerName  string
	EventName     string
	EventDate     string
	EventTime     string
	EventLocation string
	Organizer     string
	QRCode        uint
}

func SendGoMail(templatePath string, to string) {
	emailHost := os.Getenv("EMAIL")
	passwordHost := os.Getenv("PASSWORD")

	ticket := TicketDetails{
		CustomerName:  "Sakshat",
		EventName:     "Go Event Booking",
		EventDate:     "2021-09-01",
		EventTime:     "10:00",
		EventLocation: "Online",
		Organizer:     "Johnsons",
		QRCode:        1,
	}

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
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Ticket Purchase Confirmed!")
	m.SetBody("text/html", body.String())

	d := gomail.NewDialer("smtp.gmail.com", 587, emailHost, passwordHost)

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}
