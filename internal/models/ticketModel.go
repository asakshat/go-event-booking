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
	FirstName  string  `gorm:"size:50;not null" json:"first_name"`
	LastName   string  `gorm:"size:50;not null"  json:"last_name"`
	Email      string  `gorm:"size:100;not null"  json:"email"`
	Reference  string  `gorm:"size:400;not null"`
	EventPrice float64 `gorm:"type:decimal(10,2);not null"`
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

}
