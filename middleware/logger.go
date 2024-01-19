package middleware

import (
	"github.com/sing3demons/service-upload-file/logger"
	"github.com/sing3demons/service-upload-file/router"
)

func L(c router.IContext) logger.ILogger {
	switch logger := c.Get("logger").(type) {
	case logger.ILogger:
		return logger
	}
	return nil
}
