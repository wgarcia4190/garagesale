package product_test

import (
	"context"
	"github.com/wgarcia4190/garagesale/internal/schema"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/wgarcia4190/garagesale/internal/platform/database/databasetest"
	"github.com/wgarcia4190/garagesale/internal/product"
)

func TestProducts(t *testing.T) {
	db, cleanup := databasetest.Setup(t)
	defer cleanup()

	ctx := context.Background()

	np := product.NewProduct{
		Name:     "Comic Books",
		Cost:     10,
		Quantity: 20,
	}

	now := time.Date(2020, time.September, 1, 0, 0, 0, 0, time.UTC)

	p, err := product.Create(ctx, db, np, now)
	if err != nil {
		t.Fatalf("could not create product: %v", err)
	}

	saved, err := product.Retrieve(ctx, db, p.ID)
	if err != nil {
		t.Fatalf("could not retrieve product: %v", err)
	}

	if diff := cmp.Diff(p, saved); diff != "" {
		t.Fatalf("saved product did not match created: see diff:\n%s", diff)
	}
}

func TestList(t *testing.T) {
	db, cleanup := databasetest.Setup(t)
	defer cleanup()

	ctx := context.Background()

	if err := schema.Seed(db); err != nil {
		t.Fatal(err)
	}

	ps, err := product.List(ctx, db)
	if err != nil {
		t.Fatalf("listing products: %v", err)
	}

	if exp, got := 2, len(ps); exp != got {
		t.Fatalf("expected product list size %v, got %v", exp, got)
	}
}
