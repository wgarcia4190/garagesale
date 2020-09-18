package handlers

import (
	"github.com/wgarcia4190/garagesale/internal/platform/auth"
	"log"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/wgarcia4190/garagesale/internal/middleware"
	"github.com/wgarcia4190/garagesale/internal/platform/web"
)

// API constructs a handler that knows about all API routes.
func API(shutdown chan os.Signal, logger *log.Logger, db *sqlx.DB, authenticator *auth.Authenticator) http.Handler {
	app := web.NewApp(shutdown, logger, middleware.Logger(logger), middleware.Errors(logger), middleware.Metrics(),
		middleware.Panics())

	c := Check{DB: db}
	app.Handler(http.MethodGet, "/v1/health", c.Health)

	u := Users{DB: db, authenticator: authenticator}
	app.Handler(http.MethodGet, "/v1/users/token", u.Token)

	p := Product{DB: db, Log: logger}

	app.Handler(http.MethodGet, "/v1/products", p.GetListProducts, middleware.Authenticate(authenticator))
	app.Handler(http.MethodGet, "/v1/products/{id}", p.RetrieveProduct, middleware.Authenticate(authenticator))
	app.Handler(http.MethodPost, "/v1/products", p.CreateProduct, middleware.Authenticate(authenticator))
	app.Handler(http.MethodPut, "/v1/products/{id}", p.UpdateProduct, middleware.Authenticate(authenticator))
	app.Handler(http.MethodDelete, "/v1/products/{id}", p.DeleteProduct, middleware.Authenticate(authenticator),
		middleware.HasRoles(auth.RoleAdmin))

	app.Handler(http.MethodPost, "/v1/products/{id}/sales", p.AddSale, middleware.Authenticate(authenticator),
		middleware.HasRoles(auth.RoleAdmin))
	app.Handler(http.MethodGet, "/v1/products/{id}/sales", p.GetListSales, middleware.Authenticate(authenticator))

	return app
}
