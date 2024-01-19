package router

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/sing3demons/service-upload-file/files"
)

type muxContext struct {
	w http.ResponseWriter
	r *http.Request
}

func NewContext(w http.ResponseWriter, r *http.Request) IContext {
	return &muxContext{w, r}
}

func (ctx *muxContext) JSON(code int, obj any) {
	ctx.w.Header().Set("Content-Type", "application/json; charset=UTF8")
	ctx.w.WriteHeader(code)
	json.NewEncoder(ctx.w).Encode(obj)
}

func (ctx *muxContext) BodyParser(obj any) error {
	decoder := json.NewDecoder(ctx.r.Body)
	decoder.UseNumber()
	decoder.DisallowUnknownFields()
	return decoder.Decode(obj)
}

func (ctx *muxContext) ReadInput(obj any) error {
	body, err := io.ReadAll(ctx.r.Body)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(body, &obj); err != nil {
		return err
	}
	return nil
}

func (ctx *muxContext) ReadBody() ([]byte, error) {
	body, err := io.ReadAll(ctx.r.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (ctx *muxContext) Param(key string) string {
	return mux.Vars(ctx.r)[key]
}

func (ctx *muxContext) Query(key string) string {
	return ctx.r.URL.Query().Get(key)
}

func (c *muxContext) SaveUploadedFile(file multipart.File, id, path string) error {
	stor, err := files.NewLocal("./imagestore", 1024*1000*5)
	if err != nil {
		// log.Error("Unable to create storage", "error", err)
		fmt.Println("Unable to create storage", "error", err)
		os.Exit(1)
	}

	fp := filepath.Join(id, path)

	return stor.Save(fp, file)
}

func (c *muxContext) ParseMultipartForm(maxMemory int64) error {
	return c.r.ParseMultipartForm(maxMemory)
}

func (c *muxContext) Form(key string) string {
	return c.r.FormValue(key)
}

func (c *muxContext) FormFile(key string) (multipart.File, *multipart.FileHeader, error) {
	return c.r.FormFile(key)
}

