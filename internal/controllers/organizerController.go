package controllers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/asakshat/go-event-booking/initializers"
	"github.com/asakshat/go-event-booking/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func SignUp(c *gin.Context) {
	var body models.Organizer
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid body"})
		return
	}

	organizer := models.Organizer{}
	organizer.Create(initializers.DB, c, &body)

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}
func Login(c *gin.Context) {
	secret := os.Getenv("JWT_SECRET")
	var body models.Organizer
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid body"})
		return
	}
	fmt.Println(body)

	organizer := models.Organizer{}
	_, err := organizer.LoginFunc(initializers.DB, c, &body)
	if err != nil {
		return
	}
	var org models.Organizer
	initializers.DB.First(&org, "email = ?", body.Email)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": org.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Printf("Error signing token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to sign token"})
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "jwt",
		Value:    signedToken,
		MaxAge:   3600 * 24,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	})
	log.Printf("User %s logged in", body.Username)
	c.JSON(http.StatusOK, gin.H{"message": "User logged in successfully"})
}

func GetLoggedUser(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return

	}
	var org models.Organizer
	initializers.DB.First(&org, "id = ?", userID)

	c.JSON(http.StatusOK, gin.H{"user": org})
}

func Logout(c *gin.Context) {

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
		SameSite: http.SameSiteDefaultMode,
		Secure:   true,
	})

	c.JSON(http.StatusOK, gin.H{"message": "User logged out successfully"})
}
