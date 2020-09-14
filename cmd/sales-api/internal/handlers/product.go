package handlers

import (
	"github.com/pkg/errors"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
	"github.com/wgarcia4190/garagesale/internal/platform/web"
	"github.com/wgarcia4190/garagesale/internal/product"
)

// Product has handler methods for dealing with Products.
type Product struct {
	DB  *sqlx.DB
	Log *log.Logger
}

// GetListProducts gives all products as list.
func (p *Product) GetListProducts(writer http.ResponseWriter, request *http.Request) error {
	list, err := product.List(request.Context(), p.DB)

	if err != nil {
		return err
	}

	return web.Respond(writer, list, http.StatusOK)
}

// RetrieveProduct gives a single Product.
func (p *Product) RetrieveProduct(writer http.ResponseWriter, request *http.Request) error {
	id := chi.URLParam(request, "id")
	prod, err := product.Retrieve(request.Context(), p.DB, id)

	if err != nil {
		switch err {
		case product.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case product.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "looking for product %q", id)
		}
	}

	return web.Respond(writer, prod, http.StatusOK)
}

// CreateProduct decode a JSON document from a POST request and create a new Product.
func (p *Product) CreateProduct(writer http.ResponseWriter, request *http.Request) error {
	var np product.NewProduct
	if err := web.Decode(request, &np); err != nil {
		return err
	}

	prod, err := product.Create(request.Context(), p.DB, np, time.Now())
	if err != nil {
		return err
	}

	return web.Respond(writer, prod, http.StatusCreated)
}
