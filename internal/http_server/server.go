package httpserver

import (
	"fmt"
	"time"

	"anbar.bale.ai/a.iravanimanesh/concurrent-http-server/internal/api"
	"anbar.bale.ai/a.iravanimanesh/concurrent-http-server/internal/model"
	"github.com/labstack/echo"
)

func Run() {
	fmt.Println("Server is running on port " + "8080")
	time.Sleep(time.Second * 3)

	h := api.Handler{
		FileRepo: model.FileRepo{},
	}
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return h.Home(c)
	})

	e.Logger.Fatal(e.Start(":" + "8080"))
}
