package postgres

import (
	"context"
	"net"
	"net/url"
	"strconv"

	"github.com/jackc/pgx/v4"
	"github.com/tclemos/go-dockertest-example/config"
)

// NewConn creates a new postgres connection to be used across postgres repositories
func NewConn(ctx context.Context, c config.PostgresConfig) (*pgx.Conn, error) {

	dbURL := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(c.User, c.Password),
		Host:   net.JoinHostPort(c.Host, strconv.Itoa(c.Port)),
		Path:   c.Database,
	}
	q := dbURL.Query()
	q.Add("sslmode", "disable")
	dbURL.RawQuery = q.Encode()

	return pgx.Connect(ctx, dbURL.String())
}
