package postgres

import (
	"context"
	"database/sql"
	"net"
	"net/url"
	"os"

	// postgres driver
	_ "github.com/lib/pq"
	"github.com/tclemos/go-web-service-example/adapters/postgres/db"
)

// NewQuerier creates an object that contains all the query implentations
// to be consumed by the repositories
func NewQuerier(ctx context.Context) (db.Querier, error) {

	user := os.Getenv("THING_APP_POSTGRES_USER")
	password := os.Getenv("THING_APP_POSTGRES_PASSWORD")
	host := os.Getenv("THING_APP_POSTGRES_HOST")
	port := os.Getenv("THING_APP_POSTGRES_PORT")
	database := os.Getenv("THING_APP_POSTGRES_DATABASE")

	dbURL := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(user, password),
		Host:   net.JoinHostPort(host, port),
		Path:   database,
	}
	q := dbURL.Query()
	q.Add("sslmode", "disable")
	dbURL.RawQuery = q.Encode()

	d, err := sql.Open("postgres", dbURL.String())
	if err != nil {
		return nil, err
	}

	querier := db.New(d)

	return querier, nil
}
