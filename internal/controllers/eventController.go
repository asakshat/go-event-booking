package controllers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/asakshat/go-event-booking/initializers"
	"github.com/asakshat/go-event-booking/internal/models"
	"github.com/asakshat/go-event-booking/internal/services"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
)

func CreateEvent(c *gin.Context) {

	userIDUint, _ := services.GetUserID(c)

	timeString := c.Request.FormValue("time")
	dateString := c.Request.FormValue("date")

	title := c.Request.FormValue("title")
	fmt.Println("title:", title)

	fmt.Println("user:", userIDUint)

	timeStr, err := time.Parse("15:04", timeString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid time format: %v" + err.Error()})
		return
	}
	parsedTime := timeStr.Format("15:04")

	parsedDate, err := time.Parse("02/01/2006", dateString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
		return
	}

	form := models.Event{
		Title:       c.Request.FormValue("title"),
		Description: c.Request.FormValue("description"),
		Location:    c.Request.FormValue("location"),
		Venue:       c.Request.FormValue("venue"),
		Date:        parsedDate,
		Time:        parsedTime,
	}

	fileHeader, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image is required"})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer file.Close()

	tempFile, err := os.CreateTemp("", "upload-*.png")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create temp file"})
		return
	}
	defer func() {
		tempFile.Close()
		os.Remove(tempFile.Name())
	}()

	_, err = io.Copy(tempFile, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to copy file"})
		return
	}

	cld, ctx := services.Credentials()
	uniqueFilename := false
	overwrite := true
	uploadResult, err := cld.Upload.Upload(ctx, tempFile.Name(), uploader.UploadParams{
		PublicID:       form.Title,
		UniqueFilename: &uniqueFilename,
		Overwrite:      &overwrite,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image"})
		return
	}

	form.ImageURL = uploadResult.URL

	form.OrganizerID = userIDUint
	events := models.Event{}
	events.CreateEvent(initializers.DB, c, &form)

}

func EditEvent(c *gin.Context) {
	userIDUint, _ := services.GetUserID(c)

	eventID := c.Param("event_id")

	var form models.Event
	if err := initializers.DB.Where("id = ?", eventID).First(&form).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	if form.OrganizerID != userIDUint {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User is not the organizer of this event"})
		return
	}

	timeString := c.Request.FormValue("time")
	dateString := c.Request.FormValue("date")

	timeStr, err := time.Parse("15:04", timeString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid time format"})
		return
	}
	parsedTime := timeStr.Format("15:04")

	parsedDate, err := time.Parse("02/01/2006", dateString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
		return
	}

	form.Title = c.Request.FormValue("title")
	form.Description = c.Request.FormValue("description")
	form.Location = c.Request.FormValue("location")
	form.Venue = c.Request.FormValue("venue")
	form.Date = parsedDate
	form.Time = parsedTime

	fileHeader, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image is required"})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer file.Close()

	tempFile, err := os.CreateTemp("", "upload-*.png")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create temp file"})
		return
	}
	defer func() {
		tempFile.Close()
		os.Remove(tempFile.Name())
	}()

	_, err = io.Copy(tempFile, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to copy file"})
		return
	}

	cld, ctx := services.Credentials()
	uniqueFilename := false
	overwrite := true
	uploadResult, err := cld.Upload.Upload(ctx, tempFile.Name(), uploader.UploadParams{
		PublicID:       form.Title,
		UniqueFilename: &uniqueFilename,
		Overwrite:      &overwrite,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image"})
		return
	}

	form.ImageURL = uploadResult.URL

	events := models.Event{}
	events.UpdateEvent(initializers.DB, c, &form)
}

func DeleteEvent(c *gin.Context) {

	userIDUint, _ := services.GetUserID(c)

	eventID := c.Param("event_id")

	var form models.Event
	if err := initializers.DB.Where("id = ?", eventID).First(&form).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	if form.OrganizerID != userIDUint {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User is not the organizer of this event"})
		return
	}

	events := models.Event{}
	events.DeleteEvent(initializers.DB, c, &form)
}

func UndoDeleteEvent(c *gin.Context) {

	userIDUint, _ := services.GetUserID(c)

	eventID := c.Param("event_id")

	var form models.Event
	if err := initializers.DB.Unscoped().Where("id = ?", eventID).First(&form).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	if form.OrganizerID != userIDUint {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User is not the organizer of this event"})
		return
	}

	events := models.Event{}
	events.UndoDeleteEvent(initializers.DB, c, &form)
}
func DeletePerm(c *gin.Context) {
	userIDUint, _ := services.GetUserID(c)

	eventID := c.Param("event_id")

	var form models.Event
	if err := initializers.DB.Unscoped().Where("id = ?", eventID).First(&form).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	if form.OrganizerID != userIDUint {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User is not the organizer of this event"})
		return
	}

	events := models.Event{}
	events.DeletePermanently(initializers.DB, c, &form)
}
func GetEventByID(c *gin.Context) {
	eventIDStr := c.Param("event_id")
	eventIDInt, err := strconv.Atoi(eventIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}
	eventID := uint(eventIDInt)

	var event models.Event
	if err := event.GetByID(initializers.DB, eventID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"event": event})
}

func GetAllEvents(c *gin.Context) {
	var events []models.Event
	initializers.DB.Preload("Organizer").Find(&events)
	c.JSON(http.StatusOK, gin.H{"events": events})
}
