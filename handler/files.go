package handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/sing3demons/service-upload-file/logger"
	"github.com/sing3demons/service-upload-file/middleware"
	"github.com/sing3demons/service-upload-file/router"
)

// Files is a handler for reading and writing files
type Files struct{}

type File struct {
	ID      string `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Size    int64  `json:"size,omitempty"`
	IsDir   bool   `json:"is_dir,omitempty"`
	Mode    string `json:"mode,omitempty"`
	ModTime string `json:"mod_time,omitempty"`
	Href    string `json:"href,omitempty"`
}

// NewFiles creates a new File handler
func NewFiles() *Files {
	return &Files{}
}

// UploadMultipar something
func (f *Files) UploadMultipart(ctx router.IContext) {
	log := middleware.L(ctx)

	log.Info("Handle POST /upload")

	err := ctx.ParseMultipartForm(128 * 1024)
	if err != nil {
		log.Error("bad request", logger.LoggerFields{"error": err})
		// http.Error(rw, "Expected multipart form data", http.StatusBadRequest)
		ctx.JSON(http.StatusBadRequest, "Expected multipart form data")
		return
	}

	id := ctx.Form("container")

	if err != nil {
		log.Error("bad request", logger.LoggerFields{"error": err})
		ctx.JSON(http.StatusBadRequest, "invalid id, must be an integer")
		return
	}

	ff, mh, err := ctx.FormFile("file")
	if err != nil {
		log.Error("Bad request", logger.LoggerFields{"error": err})
		ctx.JSON(http.StatusBadRequest, "Expected file")
		return
	}

	ctx.SaveUploadedFile(ff, id, mh.Filename)

	ctx.JSON(http.StatusOK, fmt.Sprintf("%s/images/%s/%s", os.Getenv("HOST_URL"), id, mh.Filename))
}

func (f *Files) GetFiles(ctx router.IContext) {
	log := middleware.L(ctx)
	id := ctx.Param("container")

	pathDir, _ := os.Getwd()
	path := pathDir + "/imagestore/" + id

	// Open the directory
	dir, err := os.Open(path)
	if err != nil {
		log.Error("Unable to open directory", logger.LoggerFields{"error": err})
		return
	}
	defer dir.Close()

	// Read the entries in the directory
	fileInfos, err := dir.Readdir(0)
	if err != nil {
		log.Error("Error reading directory:", logger.LoggerFields{"error": err})
		return
	}

	var filesList []File

	// Print information about each file
	for _, fileInfo := range fileInfos {
		var file File
		file.ID = strings.Split(fileInfo.Name(), ".")[0]
		file.Name = fileInfo.Name()
		file.Size = fileInfo.Size()
		file.IsDir = fileInfo.IsDir()
		file.Mode = fileInfo.Mode().String()
		file.ModTime = fileInfo.ModTime().String()
		file.Href = fmt.Sprintf("%s/images/%s/%s", os.Getenv("HOST_URL"), id, fileInfo.Name())
		filesList = append(filesList, file)

	}

	log.Info("Handle GET", logger.LoggerFields{"container": id, "data": filesList})

	var resp []File
	for _, file := range filesList {
		resp = append(resp, File{
			ID:   file.ID,
			Name: file.Name,
			Size: file.Size,
			Href: file.Href,
		})
	}

	ctx.JSON(http.StatusOK, resp)
}

func (f *Files) GetFile(ctx router.IContext) {
	log := middleware.L(ctx)

	pwd, err := os.Getwd()
	if err != nil {
		log.Error("bad request", logger.LoggerFields{"error": err})
		ctx.JSON(400, "Invalid id, must be an integer")
		return
	}

	streamFileBytes, err := os.ReadFile(filepath.Join(pwd, "imagestore", ctx.Param("container"), ctx.Param("id")))
	if err != nil {
		log.Error("Bad request", logger.LoggerFields{"error": err})
		ctx.JSON(400, "Invalid id, must be an integer")
		return
	}
	log.Info("Handle GET", logger.LoggerFields{"container": ctx.Param("container"), "id": ctx.Param("id")})
	ctx.JSON(200, streamFileBytes)
}
