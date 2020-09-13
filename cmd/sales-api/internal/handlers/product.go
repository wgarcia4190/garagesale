package handlers

import (
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
	DB *sqlx.DB
	Log *log.Logger
}

// GetListProducts gives all products as list.
func (p *Product) GetListProducts(writer http.ResponseWriter, _ *http.Request) {
	list, err := product.List(p.DB)

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		p.Log.Println("error querying db", err)
		return
	}

	if err := web.Respond(writer, list, http.StatusOK); err != nil {
		p.Log.Println("Error writing", err)
		return
	}
}

// RetrieveProduct gives a single Product.
func (p *Product) RetrieveProduct(writer http.ResponseWriter, request *http.Request) {
	id := chi.URLParam(request, "id")
	prod, err := product.Retrieve(p.DB, id)

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		p.Log.Println("error querying db", err)
		return
	}

	if err := web.Respond(writer, prod, http.StatusOK); err != nil {
		p.Log.Println("Error writing", err)
		return
	}
}

// CreateProduct decode a JSON document from a POST request and create a new Product.
func (p *Product) CreateProduct(writer http.ResponseWriter, request *http.Request) {
	var np product.NewProduct
	if err := web.Decode(request, &np); err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		p.Log.Println("error writing", err)
	}

	prod, err := product.Create(p.DB, np, time.Now())
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		p.Log.Println("error querying db", err)
		return
	}

	if err := web.Respond(writer, prod, http.StatusCreated); err != nil {
		p.Log.Println("Error writing", err)
		return
	}
}
