package httpserver

import (
	"fmt"

	"anbar.bale.ai/a.iravanimanesh/concurrent-http-server/internal/api"
	"anbar.bale.ai/a.iravanimanesh/concurrent-http-server/internal/model"
	"github.com/labstack/echo"
)

func Run() {
	fmt.Println("Server is running on port " + "8080")

	fileRepo := model.FileRepo{
		Files: make(map[uint64]model.File),
	}

	h := api.Handler{
		FileRepo: fileRepo,
	}
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return h.Home(c)
	})

	e.POST("/upload", func(c echo.Context) error {
		return h.Upload(c)
	})

	e.POST("/download", func(c echo.Context) error {
		return h.Download(c)
	})

	e.Logger.Fatal(e.Start(":" + "8080"))
}
