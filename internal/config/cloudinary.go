package config

import (
	"log"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/spf13/viper"
)

func NewCloudinary(config *viper.Viper) *cloudinary.Cloudinary {
	cld, err := cloudinary.NewFromParams(
		config.GetString("CLOUDINARY_CLOUD_NAME"),
		config.GetString("CLOUDINARY_API_KEY"),
		config.GetString("CLOUDINARY_API_SECRET"),
	)
	if err != nil {
		log.Fatalf("❌ Failed to initialize Cloudinary: %v", err)
	}
	return cld
}
