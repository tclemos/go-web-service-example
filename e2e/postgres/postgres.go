package postgres

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"

	migrate "github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/pkg/errors"
	"github.com/sethvargo/go-retry"

	// postgres required to execute migrations
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	port = 5432
)

// Params needed to start a postgres container
type Params struct {
	Port                int
	User                string
	Password            string
	Database            string
	MigrationsDirectory string
}

// Container metadata to load a container for postgres database
type Container struct {
	name   string
	params Params
}

// NewContainer creates a new instance of Container
func NewContainer(n string, p Params) *Container {
	return &Container{
		name:   n,
		params: p,
	}
}

// Name of the container
func (c *Container) Name() string {
	return c.name
}

// Start creates and initializes a docker container for a postgres db
func (c *Container) Start(ctx context.Context, pool *dockertest.Pool) (*dockertest.Resource, error) {

	// create container resource
	resource, err := c.createResource(pool)
	if err != nil {
		return nil, err
	}

	// db url
	url := c.createDBURL(resource)

	// check db connection
	err = checkDb(ctx, url)
	if err != nil {
		return nil, err
	}

	// run migrations when a directory is specified
	if strings.TrimSpace(c.params.MigrationsDirectory) != "" {
		err = runMigrations(c.params.MigrationsDirectory, url)
		if err != nil {
			return nil, err
		}
	}

	return resource, nil
}

func (c *Container) createResource(pool *dockertest.Pool) (*dockertest.Resource, error) {

	strPort := strconv.Itoa(c.params.Port)
	pb := map[docker.Port][]docker.PortBinding{}
	pb[docker.Port(fmt.Sprintf("%d/tcp", port))] = []docker.PortBinding{
		{
			HostIP:   "0.0.0.0",
			HostPort: strPort,
		},
	}

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Name:       c.name,
		Repository: "postgres",
		Tag:        "13.2-alpine",
		Env: []string{
			"POSTGRES_DB=" + c.params.Database,
			"POSTGRES_USER=" + c.params.User,
			"POSTGRES_PASSWORD=" + c.params.Password,
		},
		PortBindings: pb,
	})

	if err != nil {
		err = errors.Wrap(err, "failed to start postgres container, check if docker is installed, running and exposing deamon on tcp://localhost:2375")
		return nil, err
	}

	return resource, nil
}

func (c *Container) createDBURL(container *dockertest.Resource) url.URL {
	// find db host
	id := fmt.Sprintf("%d/tcp", port)
	h := container.GetBoundIP(id)
	p := container.GetPort(id)
	host := net.JoinHostPort(h, p)

	// Build the connection URL.
	dbURL := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(c.params.User, c.params.Password),
		Host:   host,
		Path:   c.params.Database,
	}
	q := dbURL.Query()
	q.Add("sslmode", "disable")
	dbURL.RawQuery = q.Encode()
	return dbURL
}

func checkDb(ctx context.Context, dbURL url.URL) error {
	// prepare a connection verification interval. Use a Fibonacci backoff
	// instead of exponential so wait times scale appropriately.
	b, err := retry.NewFibonacci(500 * time.Millisecond)
	if err != nil {
		err = errors.Wrap(err, "failed to configure retries to check db connection")
		return err
	}

	b = retry.WithMaxRetries(10, b)
	b = retry.WithCappedDuration(10*time.Second, b)

	// Establish a connection to the database.
	err = retry.Do(ctx, b, func(ctx context.Context) error {
		_, err := pgxpool.Connect(ctx, dbURL.String())
		if err != nil {
			return retry.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		err = errors.Wrap(err, "failed to start postgres")
		return err
	}
	return nil
}

func runMigrations(migrationsDirectory string, dbURL url.URL) error {
	p := fmt.Sprintf("file://%s", migrationsDirectory)
	m, err := migrate.New(p, dbURL.String())
	if err != nil {
		err = errors.Wrap(err, "failed to connect to database")
		return err
	}

	err = m.Up()
	if err != nil {
		err = errors.Wrap(err, "failed to migrate database")
		return err
	}
	return nil
}
