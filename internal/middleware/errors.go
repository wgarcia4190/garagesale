package middleware

import (
	"context"
	"log"
	"net/http"

	"github.com/wgarcia4190/garagesale/internal/platform/web"
	"go.opencensus.io/trace"
)

// Errors handles errors coming out of the call chain. It detects normal
// application errors which are used to respond to the client in a uniform way.
// Unexpected errors (status >= 500) are logged.
func Errors(log *log.Logger) web.Middleware {
	// This is the actual middleware function to be executed.
	f := func(before web.Handler) web.Handler {
		h := func(ctx context.Context, writer http.ResponseWriter, request *http.Request) error {
			ctx, span := trace.StartSpan(ctx, "internal.mid.errros")
			defer span.End()

			// Run the handler chain and catch any propagated error.
			if err := before(ctx, writer, request); err != nil {
				// Log the error
				log.Printf("ERROR: %+v", err)

				// Respond to the error.
				if err := web.RespondError(ctx, writer, err); err != nil {
					return err
				}

				// Ensure that shutdown errors are allowed to bubble up to web.go
				if web.IsShutdown(err) {
					return err
				}
			}
			// Return nil to indicate the error has been handled.
			return nil
		}

		return h
	}

	return f
}
