package router

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
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
}

func NewMicroservice() IRouter {
	r := mux.NewRouter()
	mw := GzipHandler{}
	mw.GzipMiddleware(r)

	r.Use(handlers.RecoveryHandler(handlers.PrintRecoveryStack(true)))
	r.Use(mw.GzipMiddleware)
	return &router{r}
}

func (ms *router) Static(staticDir string) {
	ms.Router.PathPrefix(staticDir + "/").Handler(http.StripPrefix(staticDir+"/", http.FileServer(http.Dir("./imagestore"))))
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
		Addr:         "127.0.0.1:8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}

	go func() {
		fmt.Printf("http listen: %s\n", srv.Addr)

		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("server listen err: %v\n", err)
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("server forced to shutdown: ", err)
	}
	fmt.Println("server exited")
}
