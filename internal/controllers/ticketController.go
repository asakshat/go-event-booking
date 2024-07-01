package controllers

import (
	"fmt"
	"log"
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
	initializers.DB.Preload("Organizer").Where("id = ?", event_id).First(&event)

	var ticket models.Ticket
	if err := c.Bind(&ticket); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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
	fmt.Println(event.Organizer.Username)
	templatePath := "./html/ticket.html"
	fmt.Println(ticketDetails.Organizer)

	services.SendGoMail(templatePath, ticketDetails)

	if err := os.Remove(qrpath); err != nil {
		fmt.Printf("Error deleting QR code file: %s\n", err)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ticket bought successfully"})
}
func VerifyTicket(c *gin.Context) {
	rawToken := c.Param("token")
	ticketToken := strings.TrimPrefix(rawToken, "/")

	var ticket models.Ticket
	result := initializers.DB.Where("reference = ?", ticketToken).First(&ticket)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ticket not found"})
		return
	}

	err := services.ValidateToken(ticketToken)
	if err != nil {
		log.Println("Error validating token:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	if ticket.Validated {
		c.JSON(http.StatusConflict, gin.H{"error": "Ticket already validated"})
		return
	}
	ticket.Validated = true
	updateResult := initializers.DB.Save(&ticket)
	if updateResult.Error != nil {
		log.Println("Error updating ticket validation status:", updateResult.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating ticket validation status"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Ticket validated successfully"})
}
