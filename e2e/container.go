package e2e

import (
	"context"
	"fmt"

	"github.com/ory/dockertest/v3"
	"github.com/tclemos/go-dockertest-example/config"
	"github.com/tclemos/go-dockertest-example/logger"
)

var (
	p  *dockertest.Pool
	cc []*dockertest.Resource

	Config config.Config
	Ctx    context.Context
)

// Container represents a docker container that can be Started and Stoped
type Container interface {
	// Create and start the container
	Start(context.Context, *dockertest.Pool) (*dockertest.Resource, error)

	// Name returns the container name
	Name() string
}

func StartContainers(containers ...Container) {
	cc = []*dockertest.Resource{}
	Ctx = context.Background()

	var err error
	p, err = dockertest.NewPool("")
	if err != nil {
		panic(fmt.Errorf("failed to create Docker pool: %w", err))
	}

	for _, c := range containers {
		r, err := c.Start(Ctx, p)
		if err != nil {
			logger.Errorf(err, "failed to load container: %s :", c.Name())
			StopContainers()
			return
		}
		cc = append(cc, r)
	}
}

func StopContainers() {
	for _, c := range cc {
		err := p.Purge(c)
		if err != nil {
			logger.Errorf(err, "Could not purge container: %v", c.Container.Name)
		}
	}
}
