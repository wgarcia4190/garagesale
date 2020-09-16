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

	return web.Respond(request.Context(), writer, list, http.StatusOK)
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

	return web.Respond(request.Context(), writer, prod, http.StatusOK)
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

	return web.Respond(request.Context(), writer, prod, http.StatusCreated)
}

// UpdateProduct decodes the body of a request to update an existing product. The ID
// of the product is part of the request URL
func (p *Product) UpdateProduct(writer http.ResponseWriter, request *http.Request) error {
	id := chi.URLParam(request, "id")

	var update product.UpdateProduct
	if err := web.Decode(request, &update); err != nil {
		return errors.Wrap(err, "decoding product update")
	}

	if err := product.Update(request.Context(), p.DB, id, update, time.Now()); err != nil {
		switch err {
		case product.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case product.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "updating product %q", id)
		}
	}

	return web.Respond(request.Context(), writer, nil, http.StatusNoContent)
}

// DeleteProduct removes a single product identified by an ID in the request URL.
func (p *Product) DeleteProduct(writer http.ResponseWriter, request *http.Request) error {
	id := chi.URLParam(request, "id")

	if err := product.Delete(request.Context(), p.DB, id); err != nil {
		switch err {
		case product.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "deleting product %q", id)
		}
	}

	return web.Respond(request.Context(), writer, nil, http.StatusNoContent)
}

// AddSale creates a new Sale for a particular product. It looks for a JSON
// object in the request body. The full model is returned to the caller.
func (p *Product) AddSale(writer http.ResponseWriter, request *http.Request) error {
	var ns product.NewSale
	if err := web.Decode(request, &ns); err != nil {
		return errors.Wrap(err, "decoding new sale")
	}

	productID := chi.URLParam(request, "id")

	sale, err := product.AddSale(request.Context(), p.DB, ns, productID, time.Now())
	if err != nil {
		return errors.Wrap(err, "adding new sale")
	}

	return web.Respond(request.Context(), writer, sale, http.StatusCreated)
}

// ListSales get all sales for a particular product.
func (p *Product) GetListSales(writer http.ResponseWriter, request *http.Request) error {
	id := chi.URLParam(request, "id")

	list, err := product.ListSales(request.Context(), p.DB, id)
	if err != nil {
		return errors.Wrap(err, "getting sales list")
	}

	return web.Respond(request.Context(), writer, list, http.StatusOK)
}
