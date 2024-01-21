package main

import (
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/sing3demons/service-upload-file/handler"
	"github.com/sing3demons/service-upload-file/logger"
	"github.com/sing3demons/service-upload-file/middleware"
	"github.com/sing3demons/service-upload-file/router"
)

func init() {
	if os.Getenv("ENV_MODE") != "production" {
		if err := godotenv.Load(".env.dev"); err != nil {
			panic(err)
		}
	}
}

func main() {
	_, err := os.Create("/tmp/live")
	if err != nil {
		os.Exit(1)
	}
	defer os.Remove("/tmp/live")

	log := logger.NewLogger()
	defer log.Sync()

	r := router.NewMicroservice(log)
	h := handler.NewFiles()

	// os.MkdirAll("./imagestore", os.ModePerm)

	r.Static("/images")

	r.GET("/healthz", func(c router.IContext) {
		middleware.L(c).Info("healthz")
		c.JSON(http.StatusOK, "OK")
	})

	r.GET("/files", h.GetFolders)
	r.GET("/files/{container}", h.GetFiles)

	r.POST("/upload", h.UploadMultipart)
	r.GET("/stream/{container}/{id}", h.GetFile)

	r.StartHttp()
}
