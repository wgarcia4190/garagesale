package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/wgarcia4190/garagesale/internal/product"
)

// ProductService has handler methods for dealing with Products.
type Product struct {
	DB *sqlx.DB
}

// GetListProducts gives all products as list.
func (p *Product) GetListProducts(writer http.ResponseWriter, _ *http.Request) {
	list, err := product.List(p.DB)

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Println("error querying db", err)
		return
	}

	data, err := json.Marshal(list)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Println("Error marshalling", err)

		return
	}

	writer.Header().Set("content-type", "application/json; charset=utf-8")
	writer.WriteHeader(http.StatusOK)

	if _, err := writer.Write(data); err != nil {
		log.Println("Error writing", err)
	}
}
