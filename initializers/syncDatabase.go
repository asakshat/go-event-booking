package initializers

import "github.com/asakshat/go-event-booking/internal/models"

func SyncDatabase() {
	DB.AutoMigrate(&models.Organizer{}, models.Event{}, models.Ticket{}, models.EmailVerification{})
}
