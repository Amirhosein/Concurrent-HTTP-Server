package api

import (
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"anbar.bale.ai/a.iravanimanesh/concurrent-http-server/internal/model"
	"anbar.bale.ai/a.iravanimanesh/concurrent-http-server/internal/pkg"
	"anbar.bale.ai/a.iravanimanesh/concurrent-http-server/internal/request"
	"anbar.bale.ai/a.iravanimanesh/concurrent-http-server/internal/response"
	"github.com/labstack/echo"
)

type Handler struct {
	FileRepo model.FileRepo
}

func (h Handler) Home(c echo.Context) error {
	message := response.Message{Message: "Welcome to your concurrent http server!"}

	return c.JSON(http.StatusOK, message)
}

func (h Handler) Upload(c echo.Context) error {
	var filename string

	uploadRequest := new(request.UploadRequest)

	data, err := c.FormFile("file")
	if err == nil {
		temp, _ := data.Open()
		uploadRequest.File, _ = ioutil.ReadAll(temp)
		filename = data.Filename
	} else {
		err := c.Bind(uploadRequest)
		if err != nil {
			log.Print(err)
		}

		if err := uploadRequest.Validate(); err != nil {

			errorResponse := response.Error{
				Error: err.Error(),
			}

			return c.JSON(http.StatusNotAcceptable, errorResponse)
		}

		resp, err := http.Get(uploadRequest.File.(string))
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		filename = strings.Split(uploadRequest.File.(string), "/")[len(strings.Split(uploadRequest.File.(string), "/"))-1]

		uploadRequest.File, _ = ioutil.ReadAll(resp.Body)
	}

	file := model.File{
		Data: uploadRequest.File.([]byte),
	}

	accessHash := pkg.GenerateFileId(filename)
	h.FileRepo.Set(filename, accessHash, file)

	response := response.SuccessfulUpload{
		FileId: strconv.FormatUint(accessHash, 10) + ":" + filename,
	}

	return c.JSON(http.StatusOK, response)
}

func (h Handler) Download(c echo.Context) error {
	downloadRequest := new(request.DownloadRequest)

	err := c.Bind(downloadRequest)
	if err != nil {
		log.Print(err)
	}

	fileIdS := strings.Split(downloadRequest.FileId, ":")[0]
	fileId, _ := strconv.ParseUint(fileIdS, 10, 64)
	log.Println(fileId)

	file, ok := h.FileRepo.Get(fileId)
	if !ok {
		errorResponse := response.Error{
			Error: "File not found",
		}

		return c.JSON(http.StatusNotFound, errorResponse)
	}

	return c.Blob(http.StatusOK, "application/octet-stream", file.Data)
}
