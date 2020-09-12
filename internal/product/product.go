package product

import "github.com/jmoiron/sqlx"

// List return all known Products.
func List(db *sqlx.DB) ([]Product, error) {
	list := make([]Product, 0)

	const q = `SELECT product_id, name, cost, quantity, date_created, date_updated FROM products`

	if err := db.Select(&list, q); err != nil {
		return nil, err
	}

	return list, nil
}
