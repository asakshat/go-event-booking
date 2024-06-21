package models

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Organizer struct {
	gorm.Model
	Username     string `gorm:"size:50;not null;unique" json:"username"`
	Email        string `gorm:"size:100;not null;unique" json:"email"`
	PasswordHash string `gorm:"size:255;not null" json:"password"`
}

func (o *Organizer) Create(db *gorm.DB, c *gin.Context, body *Organizer) (*Organizer, error) {
	// check if user exists in the database
	var org Organizer
	db.First(&org, "email = ?", body.Email)
	if org.ID != 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Failed to create user. Email already exists"})
		return nil, errors.New("email already exists")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return nil, err
	}

	organizer := Organizer{Username: body.Username, Email: body.Email, PasswordHash: string(hash)}

	result := db.Create(&organizer)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create user. Username already exists"})
		return nil, result.Error
	}
	return &organizer, nil
}

func (o *Organizer) LoginFunc(db *gorm.DB, c *gin.Context, body *Organizer) (*Organizer, error) {
	var organizer Organizer
	db.First(&organizer, "email = ?", body.Email)
	if organizer.ID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return nil, errors.New("user not found")
	}
	fmt.Println(organizer)
	fmt.Println(body)
	err := bcrypt.CompareHashAndPassword([]byte(organizer.PasswordHash), []byte(body.PasswordHash))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return nil, err
	}
	return &organizer, nil
}
