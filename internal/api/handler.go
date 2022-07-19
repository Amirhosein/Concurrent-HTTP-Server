package api

import (
	"anbar.bale.ai/a.iravanimanesh/concurrent-http-server/internal/model"
	"anbar.bale.ai/a.iravanimanesh/concurrent-http-server/internal/response"
	"github.com/labstack/echo"
)

type Handler struct {
	FileDB model.FileDB
}

func (h Handler) Home(c echo.Context) error {
	message := response.Message{Message: "Hello World"}

	return c.JSON(200, message)
}
