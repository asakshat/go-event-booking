package services

import (
	"image/color"

	"github.com/skip2/go-qrcode"
)

func GenerateQRCode(data string) ([]byte, error) {
	err := qrcode.WriteColorFile("https://example.org", qrcode.Medium, 256, color.Black, color.White, "qr.png")
	if err != nil {
		return nil, err
	}
	return nil, nil
}
