package services

import (
	"github.com/skip2/go-qrcode"
)

func GenerateQRCode(token string, filename string) error {
	err := qrcode.WriteFile(token, qrcode.Medium, 256, filename+".png")
	if err != nil {
		return err
	}

	return nil
}
