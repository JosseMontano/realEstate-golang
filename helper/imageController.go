package helper

import (
	"context"
	"mime/multipart"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"

)

func Upload(formFile multipart.File) (string, bool) {
	var ctx = context.Background()

	cld, _ := cloudinary.NewFromParams("dny08tnju", "967931747444356", "4kXoesfEisSWTaSCPsq8Gno4qww")

	resp, err := cld.Upload.Upload(ctx, formFile, uploader.UploadParams{Folder: "real_estates"})

	if err != nil {
		return err.Error(), true
	}

	return resp.SecureURL, false

}
