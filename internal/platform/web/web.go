package web

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
)

// ctxKey represents the type of value for the context key.
type ctxKey int

// KeyValues is how request values or stored/retrieved.
const KeyValues ctxKey = 1

// Values carries information about each request.
type Values struct {
	StatusCode int
	Start      time.Time
}

// Handler is the signature that all application handlers will implement.
type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

// App is the entry point for all web applications.
type App struct {
	mux *chi.Mux
	Log *log.Logger
	mw  []Middleware
}

// NewApp knows how to construct internal state for an App.
func NewApp(logger *log.Logger, mw ...Middleware) *App {
	return &App{
		mux: chi.NewRouter(),
		Log: logger,
		mw:  mw,
	}
}

// Handler connects a method and URL pattern to a particular application handler.
func (a *App) Handler(method, pattern string, h Handler) {

	h = wrapMiddleware(a.mw, h)

	fn := func(w http.ResponseWriter, r *http.Request) {
		v := Values{
			Start: time.Now(),
		}

		ctx := context.WithValue(r.Context(), KeyValues, &v)

		if err := h(ctx, w, r); err != nil {
			a.Log.Printf("ERROR : Unhandled error %v", err)
		}
	}
	a.mux.MethodFunc(method, pattern, fn)
}

func (a *App) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	a.mux.ServeHTTP(writer, request)
}
