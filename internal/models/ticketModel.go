package models

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Ticket struct {
	gorm.Model
	EventID    uint    `gorm:"not null"`
	Event      Event   `gorm:"foreignKey:EventID"`
	FirstName  string  `gorm:"size:50;not null" binding:"required"`
	LastName   string  `gorm:"size:50;not null" binding:"required"`
	Email      string  `gorm:"size:100;not null" binding:"required"`
	Reference  string  `gorm:"size:255;not null"`
	EventPrice float64 `gorm:"type:decimal(10,2);not null" binding:"required"`
	QRCode     string  `gorm:"size:255"`
	Paid       bool    `gorm:"default:false"`
	Validated  bool    `gorm:"default:false"`
	BoughtDate time.Time
}
type TicketDetails struct {
	TicketID      uint
	EventName     string
	Organizer     string
	EventDate     string
	EventTime     string
	EventLocation string
	FirstName     string
	Email         string
	QRCodePath    string
}

func (t *Ticket) CreateTicket(db *gorm.DB, c *gin.Context, ticket *Ticket) {
	result := db.Create(ticket)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create ticket"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Ticket created successfully"})

}
