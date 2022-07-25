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

type LoginRegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (r LoginRegisterRequest) Validate() error {
	err := validation.Validate(r.Username, validation.Required)
	if err != nil {
		return err
	}

	err = validation.Validate(r.Password, validation.Required)
	if err != nil {
		return err
	}

	return nil
}

type AddPermissionRequest struct {
	FileId   string `json:"file_id" binding:"required"`
	Username string `json:"username_to_be_add" binding:"required"`
}

func (r AddPermissionRequest) Validate() error {
	err := validation.Validate(r.FileId, validation.Required)
	if err != nil {
		return err
	}
	err = validation.Validate(r.Username, validation.Required)
	if err != nil {
		return err
	}

	return nil
}
