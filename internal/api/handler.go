package api

import (
	"net/http"

	"anbar.bale.ai/a.iravanimanesh/concurrent-http-server/internal/model"
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
