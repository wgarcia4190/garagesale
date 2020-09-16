package middleware

import (
	"errors"
	"github.com/wgarcia4190/garagesale/internal/platform/web"
	"log"
	"net/http"
	"time"
)

// Logger will log a line for every request.
func Logger(log *log.Logger) web.Middleware {
	// This is the actual middleware function to be executed.
	f := func(before web.Handler) web.Handler {
		h := func(writer http.ResponseWriter, request *http.Request) error {

			v, ok := request.Context().Value(web.KeyValues).(*web.Values)
			if !ok {
				return errors.New("web values missing from context")
			}

			// Run the handler chain and catch any propagated error.
			err := before(writer, request)

			log.Printf(
				"%d %s %s (%v)",
				v.StatusCode, request.Method, request.URL.Path, time.Since(v.Start),
			)

			// Return the error to the handler further up the chain.
			return err
		}
		return h
	}
	return f
}
