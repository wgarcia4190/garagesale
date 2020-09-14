package databasetest

import (
	"github.com/wgarcia4190/garagesale/internal/platform/database"
	"github.com/wgarcia4190/garagesale/internal/schema"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
)

// Setup creates a test database inside a Docker container. It creates the
// required table structure but the database otherwise empty.
//
// It does not return errors as this intended for testing only. Instead it will
// call Fatal on the provided testing.T if anything goes wrong.
//
// It returns the database to user as well as a function to call at the end of the test.
func Setup(t *testing.T) (*sqlx.DB, func()) {
	t.Helper()

	c := startContainer(t)

	db, err := database.Open(database.Config{
		User:       "postgres",
		Password:   "postgres",
		Host:       c.Host,
		Name:       "postgres",
		DisableTLS: true,
	})

	if err != nil {
		t.Fatalf("opening database connection: %v", err)
	}

	t.Log("waiting for database to be ready")

	// Wait for the database to be ready. Wait 100ms longer between each attempt.
	// Do not try more than 20 times.
	var pingError error
	maxAttempts := 20

	for attempts := 1; attempts <= maxAttempts; attempts++ {
		pingError = db.Ping()
		if pingError == nil {
			break
		}
		time.Sleep(time.Duration(attempts) * 100 * time.Millisecond)
	}

	if pingError != nil {
		stopContainer(t, c)
		t.Fatalf("waiting for database to be ready: %v", err)
	}

	if err := schema.Migrate(db); err != nil {
		stopContainer(t, c)
		t.Fatalf("migrating: %s", err)
	}

	// teardown is the function that should be invoked when the caller is done
	// with the database.
	teardown := func() {
		t.Helper()
		db.Close()
		stopContainer(t, c)
	}

	return db, teardown
}
