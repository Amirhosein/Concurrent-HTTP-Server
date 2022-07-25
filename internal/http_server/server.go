package httpserver

import (
	"fmt"
	"time"

	db "anbar.bale.ai/a.iravanimanesh/concurrent-http-server/internal/DB"
	"anbar.bale.ai/a.iravanimanesh/concurrent-http-server/internal/api"
	"anbar.bale.ai/a.iravanimanesh/concurrent-http-server/internal/model"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Run() {
	fmt.Println("Server is running on port " + "8080")

	fileRepo := model.FileRepo{
		Files: make(map[uint64]model.File),
	}

	time.Sleep(time.Second * 3)

	userDB := model.UserDB{
		DB: db.InitDB(),
	}

	userCache := model.UserCache{
		DB:    userDB,
		Redis: db.InitRedis(),
	}

	h := api.Handler{
		UserRepo: userCache,
		FileRepo: fileRepo,
	}

	e := echo.New()

	// e.Use(middleware.Logger(), middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return h.Home(c)
	})

	e.POST("/register", func(c echo.Context) error {
		return h.Register(c)
	})

	e.POST("/login", func(c echo.Context) error {
		return h.Login(c)
	})

	r := e.Group("/upload")

	config := middleware.JWTConfig{
		Claims:     &model.User{},
		SigningKey: []byte("secret"),
	}
	r.Use(middleware.JWTWithConfig(config))
	r.POST("", h.Upload)
	r.POST("/a", h.Download)

	s := e.Group("/download")

	s.Use(middleware.JWTWithConfig(config))
	s.POST("", h.Download)

	v := e.Group("/addPermission")

	v.Use(middleware.JWTWithConfig(config))
	v.POST("", h.AddPermission)

	e.Logger.Fatal(e.Start(":" + "8080"))
}
