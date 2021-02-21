package e2e

import (
	"context"
	"fmt"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/pkg/errors"
	"github.com/tclemos/go-dockertest-example/logger"
)

// We need these public variables to share information betwee
// TestMain and OtherTests, if you have a better idea, tell me
var (
	pool      *dockertest.Pool
	resources []*dockertest.Resource

	Ctx context.Context
)

type ctxKey string

const (
	valuesKey = ctxKey("values")
)

// Container represents a docker container
type Container interface {
	// Options to execute the container
	Options() (*dockertest.RunOptions, error)

	// Executed after the container is started, use it to run migrations
	// copy files, etc
	AfterStart(context.Context, *dockertest.Resource) error

	// Name returns the container name
	Name() string
}

// Start the integration test environment
func Start(containers ...Container) {
	resources = []*dockertest.Resource{}
	Ctx = context.WithValue(context.Background(), valuesKey, map[string]interface{}{})

	var err error
	pool, err = dockertest.NewPool("")
	if err != nil {
		panic(fmt.Errorf("failed to create docker pool: %w", err))
	}

	for _, c := range containers {
		fmt.Printf("loading container: %s\n", c.Name())
		o, err := c.Options()
		handleContainerErr(c.Name(), "can't load run options", err)

		r, err := startContainer(Ctx, pool, o)
		handleContainerErr(c.Name(), "can't start container", err)

		err = c.AfterStart(Ctx, r)
		handleContainerErr(c.Name(), "failed to execute AfterStarted", err)

		resources = append(resources, r)
	}
}

// Stop the integration test environment
func Stop() {
	for _, r := range resources {
		err := pool.Purge(r)
		if err != nil {
			logger.Errorf(err, "Could not purge container: %v", r.Container.Name)
		}
	}
}

// AddValue allows a value to be stored during the TestMain to be used within tests
func AddValue(key string, value interface{}) {
	values := Ctx.Value(valuesKey).(map[string]interface{})
	values[key] = value
	Ctx = context.WithValue(Ctx, valuesKey, values)
}

// GetValue allows a value to be retrieved by its key
func GetValue(key string) interface{} {
	values := Ctx.Value(valuesKey).(map[string]interface{})
	return values[key]
}

// GetValues gets all stored values
func GetValues() map[string]interface{} {
	values := Ctx.Value(valuesKey).(map[string]interface{})
	return values
}

// startContainer creates and initializes a container accordingly to the provided options
func startContainer(ctx context.Context, p *dockertest.Pool, o *dockertest.RunOptions) (*dockertest.Resource, error) {
	fmt.Printf("starting container: %s\n", o.Name)
	o.Name = ""
	r, err := p.RunWithOptions(o, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		err = errors.Wrap(err, "failed to start postgres container, check if docker is installed, running and exposing deamon on tcp://localhost:2375")
		return nil, err
	}

	//err = r.Expire(50) // drop containers after 3 minutes if the got stuck
	if err != nil {
		errors.Wrap(err, "could not setup container to expire: %s")
	}

	return r, nil
}

// handleContainerErr stops the integration environment when an error is found
func handleContainerErr(m, n string, err error) {
	if err != nil {
		logger.Errorf(err, "failed to load container(%s): %s err: %s", n, m, err.Error())
		Stop()
		return
	}
}
