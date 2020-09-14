package database

import (
	"net/url"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Register the postgres database/sql driver.
)

// Config is what we require to open a database connection.
type Config struct {
	Host       string
	Name       string
	User       string
	Password   string
	DisableTLS bool
}

// Open knows how to open a database connection.
func Open(config Config) (*sqlx.DB, error) {
	q := url.Values{}
	q.Set("sslmode", "require")

	if !config.DisableTLS {
		q.Set("sslmode", "disable")
	}

	q.Set("timezone", "utc")

	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(config.User, config.Password),
		Host:     config.Host,
		Path:     config.Name,
		RawQuery: q.Encode(),
	}

	return sqlx.Open("postgres", u.String())
}
