package product

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

var (
	ErrNotFound  = errors.New("Product not found")
	ErrInvalidID = errors.New("id provided was not a valid UUID")
)

// List return all known Products.
func List(ctx context.Context, db *sqlx.DB) ([]Product, error) {
	list := make([]Product, 0)

	const q = `SELECT
			p.product_id, p.name, p.cost, p.quantity,
			COALESCE(SUM(s.quantity), 0) AS sold,
			COALESCE(SUM(s.paid), 0) AS revenue,
			p.date_created, p.date_updated
		FROM products AS p
		LEFT JOIN sales AS s On p.product_id = s.product_id
		GROUP BY p.product_id`

	if err := db.SelectContext(ctx, &list, q); err != nil {
		return nil, err
	}

	return list, nil
}

// Retrieve returns a single Product.
func Retrieve(ctx context.Context, db *sqlx.DB, id string) (*Product, error) {

	if _, err := uuid.Parse(id); err != nil {
		return nil, ErrInvalidID
	}

	var p Product

	const q = `SELECT
			p.product_id, p.name, p.cost, p.quantity,
			COALESCE(SUM(s.quantity), 0) AS sold,
			COALESCE(SUM(s.paid), 0) AS revenue,
			p.date_created, p.date_updated
		FROM products AS p
		LEFT JOIN sales AS s On p.product_id = s.product_id
		WHERE p.product_id = $1
		GROUP BY p.product_id`

	if err := db.GetContext(ctx, &p, q, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &p, nil
}

// Create makes a new Product.
func Create(ctx context.Context, db *sqlx.DB, np NewProduct, now time.Time) (*Product, error) {
	p := Product{
		ID:          uuid.New().String(),
		Name:        np.Name,
		Cost:        np.Cost,
		Quantity:    np.Quantity,
		DateCreated: now.UTC(),
		DateUpdated: now.UTC(),
	}

	const q = `INSERT INTO products
	(product_id, name, cost, quantity, date_created, date_updated)
	VALUES($1, $2, $3, $4, $5, $6)`

	if _, err := db.ExecContext(ctx, q, p.ID, p.Name, p.Cost, p.Quantity, p.DateCreated, p.DateUpdated); err != nil {
		return nil, errors.Wrapf(err, "inserting products %v", np)
	}

	return &p, nil
}
