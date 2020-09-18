package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/pkg/errors"
	"github.com/wgarcia4190/garagesale/internal/platform/auth"
	"github.com/wgarcia4190/garagesale/internal/platform/web"
	"go.opencensus.io/trace"
)

var (
	ErrForbidden = web.NewRequestError(
		errors.New("you are not authorized for that action"), http.StatusForbidden)
)

// Authenticate validates a JWT from the Authorization header.
func Authenticate(authenticator *auth.Authenticator) web.Middleware {
	// This is the actual middleware function to be executed.
	f := func(after web.Handler) web.Handler {
		// Wrap this handler around the next one provided
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			ctx, span := trace.StartSpan(ctx, "internal.mid.Authenticate")
			defer span.End()
			//  Parse the authorization header. Expected header is of
			//c the format 'Bearer <token>.'
			parts := strings.Split(r.Header.Get("Authorization"), " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				err := errors.New("expected authorization header format: Bearer <token>")
				return web.NewRequestError(err, http.StatusUnauthorized)
			}

			_, span = trace.StartSpan(ctx, "internal.auth.ParseClaims")
			claims, err := authenticator.ParseClaims(parts[1])
			if err != nil {
				return web.NewRequestError(err, http.StatusUnauthorized)
			}
			span.End()

			// Add claims to the context so they can be retrieved later.
			ctx = context.WithValue(ctx, auth.Key, claims)

			return after(ctx, w, r)
		}
		return h
	}
	return f
}

// HasRole validates that an authenticated user has at least one role from a
// specified list. This method constructs the actual function that is used.
func HasRoles(roles ...string) web.Middleware {
	// This is the actual middleware function to be executed
	f := func(after web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			ctx, span := trace.StartSpan(ctx, "internal.mid.HasRoles")
			defer span.End()

			claims, ok := ctx.Value(auth.Key).(auth.Claims)
			if !ok {
				return errors.New("claims missing from context: HasRole called without/before Authentication")
			}

			if !claims.HasRole(roles...) {
				return ErrForbidden
			}

			return after(ctx, w, r)
		}

		return h
	}

	return f
}
