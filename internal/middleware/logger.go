package middleware

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/wgarcia4190/garagesale/internal/platform/web"
	"go.opencensus.io/trace"
)

// Logger writes some information about the request to the logs in the
// format: TraceID : (200) GET /foo -> IP ADDR (latency)
func Logger(log *log.Logger) web.Middleware {
	// This is the actual middleware function to be executed.
	f := func(before web.Handler) web.Handler {
		h := func(ctx context.Context, writer http.ResponseWriter, request *http.Request) error {
			ctx, span := trace.StartSpan(ctx, "internal.mid.logger")
			defer span.End()

			v, ok := ctx.Value(web.KeyValues).(*web.Values)
			if !ok {
				return web.NewShutdownError("web values missing from context")
			}

			// Run the handler chain and catch any propagated error.
			err := before(ctx, writer, request)

			log.Printf(
				"%s : (%d) : %s %s -> %s (%s)",
				v.TraceID, v.StatusCode, request.Method, request.URL.Path, request.RemoteAddr, time.Since(v.Start),
			)

			// Return the error to the handler further up the chain.
			return err
		}
		return h
	}
	return f
}
