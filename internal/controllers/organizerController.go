package controllers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/asakshat/go-event-booking/initializers"
	"github.com/asakshat/go-event-booking/internal/models"
	"github.com/asakshat/go-event-booking/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(c *gin.Context) {
	var body models.Organizer
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid body"})
		return
	}

	var organizer models.Organizer
	updatedOrg, err := organizer.Create(initializers.DB, c, &body)
	if err != nil {
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": updatedOrg.ID,
	})
	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		log.Printf("Error signing token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to sign token"})
		return
	}
	fmt.Println("Organizer ID: ", updatedOrg.ID)
	verifyDetails := models.EmailVerification{
		OrganizerID: updatedOrg.ID,
		Email:       body.Email,
		Token:       signedToken,
		Sent:        false,
	}
	err = verifyDetails.CreateEmailVerData(initializers.DB, c, &verifyDetails)
	if err != nil {
		return
	}
	var email models.EmailVerification
	initializers.DB.First(&email, "organizer_id = ?", updatedOrg.ID)
	err = services.SendVerificationEmail(email)
	if err != nil {
		fmt.Println(err)
		return

	}
	c.JSON(http.StatusOK, gin.H{"message": "User created successfully. Please check your email and verify your account to use other features"})
}

func VerifyEmail(c *gin.Context) {
	type TokenPayload struct {
		Token string `json:"token" binding:"required"`
	}
	var token TokenPayload
	if err := c.ShouldBindJSON(&token); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}
	var verifydetails models.EmailVerification
	initializers.DB.First(&verifydetails, "token = ?", token.Token)
	if verifydetails.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No verification data found"})
		return
	}
	if token.Token != verifydetails.Token {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token is invalid"})
		return
	}

	var userdetails models.Organizer
	err := initializers.DB.First(&userdetails, "id = ?", verifydetails.OrganizerID).Error
	fmt.Println(verifydetails.OrganizerID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	userdetails.IsVerfied = true
	initializers.DB.Save(&userdetails)
	c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully"})

}

func ForgetPassword(c *gin.Context) {
	var requestBody struct {
		Email string `json:"email" binding:"required"`
	}
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	var org models.EmailVerification
	initializers.DB.First(&org, "email = ?", requestBody.Email)
	if org.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "The email does not exist in our database"})
		return
	}

	err := services.PasswordResetMail(org)
	if err != nil {
		fmt.Println(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Password reset link sent successfully"})

}

func ChangePassword(c *gin.Context) {
	var requestBody struct {
		NewPassword string `json:"new_password" binding:"required"`
		Token       string `json:"token" binding:"required"`
		Email       string `json:"email" binding:"required"`
	}
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	var org models.EmailVerification
	initializers.DB.First(&org, "email = ?", requestBody.Email)
	if org.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "The email does not exist in our database"})
		return
	}
	if requestBody.Token != org.Token {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token is invalid"})
		return
	}

	var user models.Organizer
	initializers.DB.First(&user, "email = ?", requestBody.Email)
	if user.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	err := models.PasswordValidate(requestBody.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return

	}
	hash, err := bcrypt.GenerateFromPassword([]byte(requestBody.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user.PasswordHash = string(hash)
	initializers.DB.Save(&user)
	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})

}

func Login(c *gin.Context) {
	secret := os.Getenv("JWT_SECRET")
	var body models.Organizer
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid body"})
		return
	}

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
