package services

// Import Cloudinary and other necessary libraries
//===================
import (
	"context"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
)

func Credentials() (*cloudinary.Cloudinary, context.Context) {

	cloudinaryURL := os.Getenv("CLOUDINARY_URL")
	cld, _ := cloudinary.NewFromURL(cloudinaryURL)
	cld.Config.URL.Secure = true
	ctx := context.Background()

	return cld, ctx
}
