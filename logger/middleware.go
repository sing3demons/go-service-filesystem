package logger

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

const key = "logger"

func Middleware(logger ILogger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return LogMiddleware(next, logger)
	}
}

func LogMiddleware(next http.Handler, logger ILogger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := LogParentID(r, logger)
		r = setLoggerContext(r, l)
		next.ServeHTTP(w, r)
	})
}

func LogParentID(r *http.Request, logger ILogger) ILogger {
	xParent := r.Header.Get("X-Parent-ID")
	if xParent == "" {
		xParent = uuid.NewString()
	}
	xSpan := uuid.NewString()

	return logger.With(zap.String("parent-id", xParent), zap.String("span-id", xSpan))
}

func setLoggerContext(r *http.Request, val ILogger) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), key, val))
}
