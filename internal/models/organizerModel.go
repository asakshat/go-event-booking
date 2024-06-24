package models

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
	"unicode"

	emailverifier "github.com/AfterShip/email-verifier"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	verifier = emailverifier.
		NewVerifier().
		EnableAutoUpdateDisposable()
)

type Organizer struct {
	gorm.Model
	Username     string `gorm:"size:50;not null;unique" json:"username"`
	Email        string `gorm:"size:100;not null;unique" json:"email"`
	PasswordHash string `gorm:"size:255;not null" json:"password"`
	IsVerfied    bool   `gorm:"default:false"`
}

type EmailVerification struct {
	gorm.Model
	OrganizerID uint      `gorm:"not null"`
	Email       string    `gorm:"size:100;not null;unique" json:"email"`
	Token       string    `gorm:"size:20;not null" json:"token"`
	ExpiresAt   time.Time `gorm:"not null"`
	Type        string    `gorm:"size:20;not null"`
	Sent        bool      `gorm:"default:false"`
	Organizer   Organizer `gorm:"foreignKey:OrganizerID"`
}

func EmailValidate(email string) error {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return errors.New("invalid email format")
	}
	_, domain := parts[0], parts[1]

	ret, err := verifier.Verify(email)
	if err != nil {
		fmt.Println("verify email address failed, error is: ", err)
		return err
	}

	if !ret.Syntax.Valid {
		fmt.Println("email address syntax is invalid")
		return errors.New("invalid email syntax")
	}

	if verifier.IsDisposable(domain) {
		fmt.Printf("%s is a disposable domain\n", domain)
		return errors.New("disposable email domain not allowed")
	}
	fmt.Println(ret)

	return nil

}

func PasswordValidate(password string) error {
	var errorMessages []string

	if len(password) < 8 {
		errorMessages = append(errorMessages, "be minimum 8 characters long")
	}

	containsUppercase, containsSymbol, containsNumber := false, false, false
	symbols := "~!@#$%^&*()_+{}|:\"<>?-=[]\\;',./"

	for _, r := range password {
		switch {
		case unicode.IsUpper(r):
			containsUppercase = true
		case unicode.IsNumber(r):
			containsNumber = true
		case strings.ContainsRune(symbols, r):
			containsSymbol = true
		}
	}

	if !containsUppercase {
		errorMessages = append(errorMessages, "contain minimum 1 uppercase letter")
	}
	if !containsSymbol {
		errorMessages = append(errorMessages, "contain minimum 1 symbol")
	}
	if !containsNumber {
		errorMessages = append(errorMessages, "contain minimum 1 number")
	}

	if len(errorMessages) > 0 {
		errorMessage := "Password should " + strings.Join(errorMessages[:len(errorMessages)-1], ", ")
		if len(errorMessages) > 1 {
			errorMessage += ", and " + errorMessages[len(errorMessages)-1]
		} else {
			errorMessage += errorMessages[0]
		}
		return errors.New(errorMessage)
	}

	return nil
}

func UsernameValidate(username string) error {
	var errorMessages []string

	if len(username) < 5 {
		errorMessages = append(errorMessages, "be minimum 5 characters long")
	}

	for _, r := range username {
		if !unicode.IsLetter(r) && !unicode.IsNumber(r) {
			errorMessages = append(errorMessages, "contain only alphanumeric characters")
			break
		}
	}

	if len(errorMessages) > 0 {
		errorMessage := "Username should " + strings.Join(errorMessages[:len(errorMessages)-1], ", ")
		if len(errorMessages) > 1 {
			errorMessage += ", and " + errorMessages[len(errorMessages)-1]
		} else {
			errorMessage += errorMessages[0]
		}
		return errors.New(errorMessage)
	}

	return nil
}

func (o *Organizer) Create(db *gorm.DB, c *gin.Context, body *Organizer) (*Organizer, error) {
	var orgByUsername Organizer
	db.First(&orgByUsername, "username = ?", body.Username)
	if orgByUsername.ID != 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Failed to create user. Username already exists"})
		return nil, errors.New("username already exists")
	}

	var orgByEmail Organizer
	db.First(&orgByEmail, "email = ?", body.Email)
	if orgByEmail.ID != 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Failed to create user. Email already exists"})
		return nil, errors.New("email already exists")
	}

	err := UsernameValidate(body.Username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return nil, err
	}

	err = EmailValidate(body.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return nil, err

	}

	err = PasswordValidate(body.PasswordHash)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return nil, err
	}

	organizer := Organizer{Username: body.Username, Email: body.Email, PasswordHash: string(hash)}

	result := db.Create(&organizer)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create user"})
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
