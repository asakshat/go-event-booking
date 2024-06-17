package models

import (
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Event struct {
	gorm.Model
	OrganizerID   uint      `gorm:"not null"`
	Organizer     Organizer `gorm:"foreignKey:OrganizerID"`
	Title         string    `gorm:"size:100;not null"`
	Description   string    `gorm:"type:text;not null"`
	ImageURL      string    `gorm:"size:255"`
	Location      string    `gorm:"type:text"`
	Venue         string    `gorm:"size:100"`
	Date          time.Time `gorm:"type:date;not null"`
	Time          string    `gorm:"size:5;not null"`
	TicketDetails []Ticket
}

func (e *Event) CreateEvent(db *gorm.DB, c *gin.Context, event *Event) {
	result := db.Create(event)
	if result.Error != nil {
		c.JSON(400, gin.H{"error": "Failed to create event"})
		return
	}
	c.JSON(200, gin.H{"message": "Event created successfully"})

}

func (e *Event) UpdateEvent(db *gorm.DB, c *gin.Context, event *Event) {
	result := db.Save(event)
	if result.Error != nil {
		c.JSON(400, gin.H{"error": "Failed to update event"})
		return
	}
	c.JSON(200, gin.H{"message": "Event updated successfully"})

}

func (e *Event) DeleteEvent(db *gorm.DB, c *gin.Context, event *Event) {
	result := db.Delete(event)
	if result.Error != nil {
		c.JSON(400, gin.H{"error": "Failed to delete event"})
		return
	}
	c.JSON(200, gin.H{"message": "Event deleted successfully"})

}
func (e *Event) UndoDeleteEvent(db *gorm.DB, c *gin.Context, event *Event) {
	result := db.Model(event).Unscoped().Update("deleted_at", nil)
	if result.Error != nil {
		c.JSON(400, gin.H{"error": "Failed to undo delete event"})
		return
	}
	c.JSON(200, gin.H{"message": "Event delete state changed successfully"})
}

func (e *Event) DeletePermanently(db *gorm.DB, c *gin.Context, event *Event) {
	result := db.Unscoped().Delete(event)
	if result.Error != nil {
		c.JSON(400, gin.H{"error": "Failed to delete event"})
		return
	}
	c.JSON(200, gin.H{"message": "Event deleted permanently "})

}

func (e *Event) GetByID(db *gorm.DB, id uint) error {
	return db.Where("id = ?", id).Preload("Organizer").First(e).Error
}
