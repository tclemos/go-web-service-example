package postgres

import (
	"context"
	"database/sql"
	"net"
	"net/url"
	"strconv"

	_ "github.com/lib/pq"
	"github.com/tclemos/go-web-service-example/adapters/postgres/db"
	"github.com/tclemos/go-web-service-example/config"
)

// NewConn creates a new postgres connection to be used across postgres repositories
func NewQuerier(ctx context.Context, c config.PostgresConfig) (db.Querier, error) {

	dbURL := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(c.User, c.Password),
		Host:   net.JoinHostPort(c.Host, strconv.Itoa(c.Port)),
		Path:   c.Database,
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
