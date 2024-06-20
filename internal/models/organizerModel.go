package models

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Organizer struct {
	gorm.Model
	Username     string `gorm:"size:50;not null;unique" `
	Email        string `gorm:"size:100;not null;unique" binding:"required,email"`
	PasswordHash string `gorm:"size:255;not null" json:"-" binding:"required,min=8"`
}

func (o *Organizer) Create(db *gorm.DB, c *gin.Context, body *Organizer) {
	// check if user exists in the database
	var org Organizer
	db.First(&org, "email = ?", body.Email)
	if org.ID != 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	organizer := Organizer{Username: body.Username, Email: body.Email, PasswordHash: string(hash)}

	result := db.Create(&organizer)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create user"})
		return
	}
}

func (o *Organizer) LoginFunc(db *gorm.DB, c *gin.Context, body *Organizer) (*Organizer, error) {
	organizer := Organizer{Email: body.Email, PasswordHash: body.PasswordHash}
	result := db.Where("email = $1 ", body.Email).First(&organizer)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User or Email  not found"})
		return nil, result.Error
	}

	err := bcrypt.CompareHashAndPassword([]byte(organizer.PasswordHash), []byte(body.PasswordHash))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return nil, result.Error
	}
	return nil, result.Error
}
func (o *Organizer) GetLoggedUser(c *gin.Context) {
	// Get the user from the context
	user, _ := c.Get("user")
	c.JSON(http.StatusOK, gin.H{"user": user})
}
