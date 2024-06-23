package controllers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/asakshat/go-event-booking/initializers"
	"github.com/asakshat/go-event-booking/internal/models"
	"github.com/asakshat/go-event-booking/internal/services"
	"github.com/gin-gonic/gin"
)

func BuyTicket(c *gin.Context) {
	event_id := c.Params.ByName("event_id")

	var event models.Event
	initializers.DB.Where("id = ?", event_id).First(&event)

	var ticket models.Ticket
	if err := c.Bind(&ticket); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(ticket)

	token, err := services.GenerateToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	ticket.EventID = event.ID
	ticket.EventPrice = event.Price
	ticket.Reference = token

	filename := fmt.Sprintf("qr-code-%s-%s-%s",
		strings.ReplaceAll(ticket.FirstName, " ", "_"),
		strings.ReplaceAll(ticket.LastName, " ", "_"),
		strings.ReplaceAll(event.Title, " ", "_"))

	services.GenerateQRCode(token, filename)

	ticketAdd := models.Ticket{}
	ticketAdd.CreateTicket(initializers.DB, c, &ticket)

	dir := "./"
	qrpath := filepath.Join(dir, filename+".png")

	formattedDate := event.Date.Format("02/01/2006")
	ticketDetails := models.TicketDetails{
		TicketID:      ticketAdd.ID,
		EventName:     event.Title,
		Organizer:     event.Organizer.Username,
		EventDate:     formattedDate,
		EventTime:     event.Time,
		FirstName:     ticket.FirstName,
		EventTitle:    event.Title,
		LastName:      ticket.LastName,
		EventLocation: event.Location,
		Email:         ticket.Email,
		QRCodePath:    qrpath,
	}

	templatePath := "./index.html"
	services.SendGoMail(templatePath, ticketDetails)

	if err := os.Remove(qrpath); err != nil {
		fmt.Printf("Error deleting QR code file: %s\n", err)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ticket bought successfully"})
}
