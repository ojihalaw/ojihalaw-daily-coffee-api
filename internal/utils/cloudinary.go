package utils

import (
	"context"
	"fmt"
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

// UploadImageIntoCloudinary uploads an image to Cloudinary and returns the secure URL
func UploadImageIntoCloudinary(
	cld *cloudinary.Cloudinary,
	ctx context.Context,
	file *multipart.FileHeader,
	folder string, // folder bisa dikirim biar lebih fleksibel
) (string, error) {
	// Open the file
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	// Upload the file to Cloudinary
	resp, err := cld.Upload.Upload(ctx, src, uploader.UploadParams{
		Folder: folder,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload image: %w", err)
	}

	return resp.SecureURL, nil
}
