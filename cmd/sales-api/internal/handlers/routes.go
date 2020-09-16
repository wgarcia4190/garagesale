package handlers

import (
	"github.com/wgarcia4190/garagesale/internal/platform/auth"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/wgarcia4190/garagesale/internal/middleware"
	"github.com/wgarcia4190/garagesale/internal/platform/web"
)

// API constructs a handler that knows about all API routes.
func API(logger *log.Logger, db *sqlx.DB, authenticator *auth.Authenticator) http.Handler {
	app := web.NewApp(logger, middleware.Logger(logger), middleware.Errors(logger), middleware.Metrics())

	c := Check{DB: db}
	app.Handler(http.MethodGet, "/v1/health", c.Health)

	u := Users{DB: db, authenticator: authenticator}
	app.Handler(http.MethodGet, "/v1/users/token", u.Token)

	p := Product{DB: db, Log: logger}

	app.Handler(http.MethodGet, "/v1/products", p.GetListProducts)
	app.Handler(http.MethodGet, "/v1/products/{id}", p.RetrieveProduct)
	app.Handler(http.MethodPost, "/v1/products", p.CreateProduct)
	app.Handler(http.MethodPut, "/v1/products/{id}", p.UpdateProduct)
	app.Handler(http.MethodDelete, "/v1/products/{id}", p.DeleteProduct)

	app.Handler(http.MethodPost, "/v1/products/{id}/sales", p.AddSale)
	app.Handler(http.MethodGet, "/v1/products/{id}/sales", p.GetListSales)

	return app
}
