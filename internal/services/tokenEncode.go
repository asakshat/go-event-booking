package services

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/asakshat/go-event-booking/initializers"
	"github.com/google/uuid"
)

func GenerateToken() (string, error) {
	token := uuid.New().String()
	hashed := sha256.Sum256([]byte(token))
	signature, err := rsa.SignPSS(rand.Reader, initializers.PrivateKey, crypto.SHA256, hashed[:], nil)
	if err != nil {
		return "", err
	}
	signedToken := base64.StdEncoding.EncodeToString(signature) + "." + token
	return signedToken, nil
}

func ValidateToken(signedToken string) (bool, error) {
	parts := strings.Split(signedToken, ".")
	if len(parts) != 2 {
		return false, fmt.Errorf("invalid token format")
	}
	signature, err := base64.StdEncoding.DecodeString(parts[0])
	if err != nil {
		return false, err
	}
	hashed := sha256.Sum256([]byte(parts[1]))
	err = rsa.VerifyPSS(initializers.PublicKey, crypto.SHA256, hashed[:], signature, nil)
	if err != nil {
		return false, err
	}
	return true, nil
}
