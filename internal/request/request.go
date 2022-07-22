package request

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type UploadRequest struct {
	File interface{} `json:"file" binding:"required"`
}

func (r UploadRequest) Validate() error {
	err := validation.Validate(r.File.(string), validation.Required, is.URL)
	if err != nil {
		return err
	}
	return nil
}

type DownloadRequest struct {
	FileId string `json:"file_id" binding:"required"`
}

func (r DownloadRequest) Validate() error {
	err := validation.Validate(r.FileId, validation.Required)
	if err != nil {
		return err
	}
	return nil
}
