package lib

import (
	"context"
	"mime/multipart"
	"os"

	appErrors "github.com/buildyow/byow-user-service/domain/errors"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func CloudinaryUpload(file multipart.File) (string, error) {
	cld, err := cloudinary.NewFromParams(
		os.Getenv("CLOUDINARY_CLOUD_NAME"),
		os.Getenv("CLOUDINARY_API_KEY"),
		os.Getenv("CLOUDINARY_API_SECRET"),
	)
	if err != nil {
		return "", appErrors.WrapError(err, "Failed to initialize Cloudinary")
	}

	uploadResp, err := cld.Upload.Upload(context.Background(), file, uploader.UploadParams{})
	if err != nil {
		return "", appErrors.ErrCloudinaryUploadFailed
	}

	return uploadResp.SecureURL, nil
}
