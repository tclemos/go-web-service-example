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
	fmt.Printf("[e2e]initializing containers\n")
	resources = []*dockertest.Resource{}
	Ctx = context.WithValue(context.Background(), valuesKey, map[string]interface{}{})

	var err error
	pool, err = dockertest.NewPool("")
	if err != nil {
		panic(fmt.Errorf("[e2e]failed to create docker pool: %w", err))
	}

	for _, c := range containers {
		o, err := c.Options()
		fmt.Printf("[e2e]loading container with options: %s\n%v\n", c.Name(), o)
		handleContainerErr(c.Name(), "[e2e]can't load run options", err)

		r, err := startContainer(Ctx, pool, o)
		handleContainerErr(c.Name(), "[e2e]can't start container", err)

		fmt.Printf("[e2e]preparing container: %s\n", c.Name())
		err = c.AfterStart(Ctx, r)
		handleContainerErr(c.Name(), "[e2e]failed to execute AfterStarted", err)

		resources = append(resources, r)
	}
}

// Stop the integration test environment
func Stop() {
	for _, r := range resources {
		fmt.Printf("[e2e]purging container: %s\n", r.Container.Name)
		err := pool.Purge(r)
		if err != nil {
			logger.Errorf(err, "[e2e]could not purge container: %v", r.Container.Name)
		} else {
			fmt.Printf("[e2e]container purged: %s\n", r.Container.Name)
		}
	}
}

// startContainer creates and initializes a container accordingly to the provided options
func startContainer(ctx context.Context, p *dockertest.Pool, o *dockertest.RunOptions) (*dockertest.Resource, error) {
	fmt.Printf("[e2e]starting container: %s\n", o.Name)
	r, err := p.RunWithOptions(o, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		err = errors.Wrap(err, "[e2e]failed to start postgres container, check if docker is installed, running and exposing deamon on tcp://localhost:2375")
		return nil, err
	}

	err = r.Expire(60)
	if err != nil {
		errors.Wrap(err, "could not setup container to expire: %s")
	}

	fmt.Printf("[e2e]container started: %s\n", o.Name)
	return r, nil
}

// handleContainerErr stops the integration environment when an error is found
func handleContainerErr(m, n string, err error) {
	if err != nil {
		logger.Errorf(err, "[e2e]failed to load container(%s): %s err: %s", n, m, err.Error())
		Stop()
		return
	}
}
