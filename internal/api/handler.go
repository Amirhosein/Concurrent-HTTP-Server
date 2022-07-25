package api

import (
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"anbar.bale.ai/a.iravanimanesh/concurrent-http-server/internal/model"
	"anbar.bale.ai/a.iravanimanesh/concurrent-http-server/internal/pkg"
	"anbar.bale.ai/a.iravanimanesh/concurrent-http-server/internal/request"
	"anbar.bale.ai/a.iravanimanesh/concurrent-http-server/internal/response"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	FileRepo model.FileRepo
	UserRepo model.UserRepo
}

func (h Handler) Home(c echo.Context) error {
	message := response.Message{Message: "Welcome to your concurrent http server!"}

	return c.JSON(http.StatusOK, message)
}

func (h Handler) Register(c echo.Context) error {
	request := new(request.LoginRegisterRequest)

	err := c.Bind(request)
	if err != nil {
		log.Print(err)
	}

	if err := request.Validate(); err != nil {
		errorResponse := response.Error{
			Error: err.Error(),
		}

		return c.JSON(http.StatusNotAcceptable, errorResponse)
	}

	_, err = h.UserRepo.Set(request.Username, request.Password)
	if err != nil {
		errorResponse := response.Error{
			Error: err.Error(),
		}

		return c.JSON(http.StatusNotAcceptable, errorResponse)
	}

	response := response.Message{
		Message: "User registered successfully",
	}

	return c.JSON(http.StatusOK, response)
}

func (h Handler) Login(c echo.Context) error {
	loginRequest := new(request.LoginRegisterRequest)

	err := c.Bind(loginRequest)
	if err != nil {
		log.Print(err)
	}

	if err := loginRequest.Validate(); err != nil {
		errorResponse := response.Error{
			Error: err.Error(),
		}

		return c.JSON(http.StatusNotAcceptable, errorResponse)
	}

	user, err := h.UserRepo.Get(loginRequest.Username, loginRequest.Password)
	if err != nil {
		errorResponse := response.Error{
			Error: "User not found",
		}
		return c.JSON(http.StatusNotFound, errorResponse)
	}

	claim := &model.User{
		user.ID,
		user.Username,
		user.Files,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return err
	}

	response := response.SuccessfulLogin{
		Message: "Login successful",
		Token:   t,
	}

	return c.JSON(http.StatusOK, response)
}

func (h Handler) Upload(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*model.User)

	usr, err := h.UserRepo.Get(claims.Username, "")
	if err != nil {
		errorResponse := response.Error{
			Error: err.Error(),
		}

		return c.JSON(http.StatusNotFound, errorResponse)
	}

	files := usr.Files

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

	files = append(files, strconv.FormatUint(accessHash, 10)+":"+filename)
	claims.Files = files
	log.Println(claims.Files)

	err = h.UserRepo.Update(*claims)
	if err != nil {
		errorResponse := response.Error{
			Error: err.Error(),
		}

		return c.JSON(http.StatusNotAcceptable, errorResponse)
	}

	response := response.SuccessfulUpload{
		FileId: strconv.FormatUint(accessHash, 10) + ":" + filename,
	}

	return c.JSON(http.StatusOK, response)
}

func (h Handler) Download(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*model.User)
	username := claims.Username
	log.Println(username)
	downloadRequest := new(request.DownloadRequest)

	err := c.Bind(downloadRequest)
	if err != nil {
		log.Print(err)
	}

	usr, err := h.UserRepo.Get(username, "")
	if err != nil {
		errorResponse := response.Error{
			Error: err.Error(),
		}

		return c.JSON(http.StatusNotAcceptable, errorResponse)
	}
	files := usr.Files
	log.Println(files)

	fileIdS := strings.Split(downloadRequest.FileId, ":")[0]
	fileId, _ := strconv.ParseUint(fileIdS, 10, 64)
	log.Println(fileId)

	// check if fileId exists in user's files
	if !pkg.Contains(files, downloadRequest.FileId) {
		errorResponse := response.Error{
			Error: "File not found",
		}

		return c.JSON(http.StatusNotFound, errorResponse)
	}

	file, ok := h.FileRepo.Get(fileId)
	if !ok {
		errorResponse := response.Error{
			Error: "File not found",
		}

		return c.JSON(http.StatusNotFound, errorResponse)
	}

	return c.Blob(http.StatusOK, "application/octet-stream", file.Data)
}

func (h Handler) AddPermission(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*model.User)
	username := claims.Username
	log.Println(username)
	permissionRequest := new(request.AddPermissionRequest)
	err := c.Bind(permissionRequest)
	if err != nil {
		log.Print(err)
	}
	usr, err := h.UserRepo.Get(username, "")
	if err != nil {
		errorResponse := response.Error{
			Error: err.Error(),
		}

		return c.JSON(http.StatusNotAcceptable, errorResponse)
	}
	files := usr.Files
	log.Println(files)
	fileIdS := strings.Split(permissionRequest.FileId, ":")[0]
	fileId, _ := strconv.ParseUint(fileIdS, 10, 64)
	log.Println(fileId)
	// check if fileId exists in user's files
	if !pkg.Contains(files, permissionRequest.FileId) {
		errorResponse := response.Error{
			Error: "You don't have permission to this file",
		}

		return c.JSON(http.StatusNotFound, errorResponse)
	}

	secondUser, err := h.UserRepo.Get(permissionRequest.Username, "")
	if err != nil {
		errorResponse := response.Error{
			Error: err.Error(),
		}

		return c.JSON(http.StatusNotAcceptable, errorResponse)
	}

	secondUser.Files = append(secondUser.Files, permissionRequest.FileId)
	err = h.UserRepo.Update(secondUser)
	if err != nil {
		errorResponse := response.Error{
			Error: err.Error(),
		}

		return c.JSON(http.StatusNotAcceptable, errorResponse)
	}

	response := response.Message{
		Message: "Permission added successfully",
	}
	return c.JSON(http.StatusOK, response)
}
