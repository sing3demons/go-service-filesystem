package router

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sing3demons/service-upload-file/logger"
)

type IRouter interface {
	GET(path string, h handlerFunc)
	POST(path string, h handlerFunc)
	PUT(path string, h handlerFunc)
	PATCH(path string, h handlerFunc)
	DELETE(path string, h handlerFunc)
	Use(mwf ...mux.MiddlewareFunc)
	Static(path string)

	StartHttp()
}

type router struct {
	*mux.Router
	Log logger.ILogger
}

func NewMicroservice(logg logger.ILogger) IRouter {
	r := mux.NewRouter()
	mw := GzipHandler{}
	mw.GzipMiddleware(r)

	allowedOrigins := os.Getenv("ORIGIN_ALLOWED")

	r.Use(handlers.RecoveryHandler(handlers.PrintRecoveryStack(true)))
	r.Use(mw.GzipMiddleware)
	r.Use(logger.Middleware(logg))
	if allowedOrigins != "" {
		headersOk := handlers.AllowedHeaders([]string{"X-Requested-With"})
		originsOk := handlers.AllowedOrigins(strings.Split(allowedOrigins, ","))
		methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
		r.Use(handlers.CORS(originsOk, headersOk, methodsOk))
	}

	return &router{r, logg}
}

func (ms *router) Static(staticDir string) {
	handlerServe := http.StripPrefix(staticDir+"/", http.FileServer(http.Dir("./imagestore")))
	ms.Router.PathPrefix(staticDir + "/").Handler(handlerServe)
}

func (ms *router) GET(path string, h handlerFunc) {
	ms.Router.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		h(NewContext(w, r))
	}).Methods(http.MethodGet)
}

func (ms *router) Use(mwf ...mux.MiddlewareFunc) {
	ms.Router.Use(mwf...)
}

func (ms *router) POST(path string, h handlerFunc) {
	ms.Router.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		h(NewContext(w, r))
	}).Methods(http.MethodPost)
}

func (ms *router) PUT(path string, h handlerFunc) {
	ms.Router.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		h(NewContext(w, r))
	}).Methods(http.MethodPut)
}

func (ms *router) PATCH(path string, h handlerFunc) {
	ms.Router.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		h(NewContext(w, r))
	}).Methods(http.MethodPatch)
}

func (ms *router) DELETE(path string, h handlerFunc) {
	ms.Router.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		h(NewContext(w, r))
	}).Methods(http.MethodDelete)
}

func (ms *router) StartHttp() {
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	srv := &http.Server{
		Handler:      ms.Router,
		Addr:         ":" + os.Getenv("PORT"),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}

	currentTime := time.Now()
	currentTimeZone, offset := currentTime.Zone()
	loc, _ := time.LoadLocation("Asia/Bangkok")

	go func() {
		hostName, _ := os.Hostname()
		ms.Log.Info("HTTP server is running", logger.LoggerFields{
			"HOST":           hostName,
			"PORT":           srv.Addr,
			"ENV":            os.Getenv("ENV_MODE"),
			"LVL":            os.Getenv("LOG_LEVEL"),
			"TZ::OFFSET":     fmt.Sprintf("%s, %d", currentTimeZone, offset),
			"TIME::LOCATION": fmt.Sprintf("%s(%s)", currentTime.In(loc), loc.String()),
		})

		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			ms.Log.Error("listen", logger.LoggerFields{
				"error": err.Error(),
			})
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ms.Log.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		ms.Log.Error("server forced to shutdown: ", logger.LoggerFields{
			"error": err,
		})
		os.Exit(1)
	}
	ms.Log.Info("server exiting")
}
