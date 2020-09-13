package web

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

// App is the entry point for all web applications.
type App struct {
	mux *chi.Mux
	Log *log.Logger
}

// NewApp knows how to construct internal state for an App.
func NewApp(logger *log.Logger) *App {
	return &App{
		mux: chi.NewRouter(),
		Log: logger,
	}
}

// Handler connects a method and URL pattern to a particular application handler.
func (a *App) Handler(method, pattern string, fn http.HandlerFunc) {
	a.mux.MethodFunc(method, pattern, fn)
}

func (a *App) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	a.mux.ServeHTTP(writer, request)
}
