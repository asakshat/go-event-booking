package controllers

import (
	"fmt"
	"net/http"

	"github.com/asakshat/go-event-booking/initializers"
	"github.com/asakshat/go-event-booking/internal/models"
	"github.com/gin-gonic/gin"
)

// func generateAndSendQRCode(c *gin.Context) {
// 	// Generate a signed token

// 	err = services.GenerateQRCode(signedToken, "./qr-code.png")
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to generate QR code"})
// 		return
// 	}

// 	qrCode := "./qr-code.png"
// 	to := "recipient@example.com" // Replace with the recipient's email
// 	subject := "Your QR Code"
// 	body := "Please find attached your QR code."
// 	err = services.SendEmail(to, subject, body, qrCode)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to send email"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "QR code sent successfully"})
// }

type EmailDetails struct {
}

func BuyTicket(c *gin.Context) {
	var ticket models.Ticket
	if err := c.ShouldBindJSON(&ticket); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ticketAdd := models.Ticket{}
	ticketAdd.CreateTicket(initializers.DB, c, &ticket)

	var results []models.TicketDetails

	initializers.DB.Raw(`
    SELECT 
        tickets.id,
        organizers.username,
        events.title, 
        events.date, 
        events.time, 
        events.location, 
		,
		
    FROM 
        events
    INNER JOIN 
        organizers ON events.organizer_id = organizers.id
    INNER JOIN 
        tickets ON events.id = tickets.event_id
    WHERE 
        events.id = ? AND tickets.id = ?`, ticket.EventID, ticket.ID).Scan(&results)

	for _, result := range results {
		fmt.Printf("Ticket ID: %d, Event Name: %s, Organizer Name: %s, Event Date: %s, Event Time: %s, Event Location: %s , First Name: %s, Email: %s, QR Code Path: %s\n",
			result.TicketID, result.EventName, result.Organizer, result.EventDate, result.EventTime, result.EventLocation, result.FirstName, result.Email, result.QRCodePath)
	}
}
