package main

import (
	"github.com/sing3demons/service-upload-file/handler"
	"github.com/sing3demons/service-upload-file/logger"
	"github.com/sing3demons/service-upload-file/router"
)

func main() {

	log := logger.NewLogger()
	defer log.Sync()

	r := router.NewMicroservice(log)
	h := handler.NewFiles()

	r.Static("/images")

	r.GET("/files/{container}", h.GetFiles)

	r.POST("/upload", h.UploadMultipart)

	r.StartHttp()
}
