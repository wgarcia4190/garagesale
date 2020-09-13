package handlers

import (
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/wgarcia4190/garagesale/internal/platform/web"
)

// API constructs a handler that knows about all API routes.
func API(logger *log.Logger, db *sqlx.DB) http.Handler {
	app := web.NewApp(logger)
	p := Product{DB: db, Log: logger}

	app.Handler(http.MethodGet, "/v1/products", p.GetListProducts)
	app.Handler(http.MethodGet, "/v1/products/{id}", p.RetrieveProduct)
	app.Handler(http.MethodPost, "/v1/products", p.CreateProduct)

	return app
}
