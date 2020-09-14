package web

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

// Handler is the signature that all application handlers will implement.
type Handler func(w http.ResponseWriter, r *http.Request) error

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
func (a *App) Handler(method, pattern string, h Handler) {

	fn := func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			a.Log.Printf("ERROR: %v", err)

			if err := RespondError(w, err); err != nil {
				a.Log.Printf("ERROR: %v", err)
			}
		}
	}
	a.mux.MethodFunc(method, pattern, fn)
}

func (a *App) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	a.mux.ServeHTTP(writer, request)
}
