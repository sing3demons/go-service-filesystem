package router

import "mime/multipart"

type IContext interface {
	JSON(code int, obj any)
	Param(key string) string
	Query(key string) string
	Form(key string) string
	FormFile(key string) (multipart.File, *multipart.FileHeader, error)

	BodyParser(obj any) error
	ReadBody() ([]byte, error)
	SaveUploadedFile(file multipart.File, dst, path string) error
	ParseMultipartForm(maxMemory int64) error
}

type handlerFunc func(c IContext)
