package middleware

import (
	"github.com/sing3demons/service-upload-file/logger"
	"github.com/sing3demons/service-upload-file/router"
)

func L(c router.IContext) logger.ILogger {
	switch logg := c.Get(logger.Key).(type) {
	case logger.ILogger:
		return logg
	default:
		return logger.NewLogger()
	}
}
