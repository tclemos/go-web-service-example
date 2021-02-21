package localstack

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/pkg/errors"
	"github.com/tclemos/go-dockertest-example/e2e"
)

const (
	edgePort        = 4566
	internalWebPort = 8080
	externalWebPort = 8123
)

type SqsQueue struct {
	Name string
}

// Params needed to start a postgres container
type Params struct {
	Region    string
	SqsQueues []SqsQueue
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

// Options to start a localstack container accordingly to the params
func (c *Container) Options() (*dockertest.RunOptions, error) {
	pb := map[docker.Port][]docker.PortBinding{}
	pb[docker.Port(fmt.Sprintf("%d/tcp", edgePort))] = []docker.PortBinding{{HostIP: "0.0.0.0", HostPort: strconv.Itoa(edgePort)}}
	pb[docker.Port(fmt.Sprintf("%d/tcp", internalWebPort))] = []docker.PortBinding{{HostIP: "0.0.0.0", HostPort: strconv.Itoa(externalWebPort)}}
	return &dockertest.RunOptions{
		Name:       c.name,
		Repository: "localstack/localstack",
		Tag:        "latest",
		Env: []string{
			"SERVICES=sqs",
			"DATA_DIR=/tmp/localstack/data",
		},
		PortBindings: pb,
	}, nil
}

// AfterStart will check the connection and execute migrations
func (c *Container) AfterStart(ctx context.Context, r *dockertest.Resource) error {
	// create sqs queues
	for _, q := range c.params.SqsQueues {
		if _, err := r.Exec([]string{"bash", "-c", fmt.Sprintf("awslocal sqs create-queue --queue-name %s", q.Name)}, dockertest.ExecOptions{}); err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to write config to sqs container: %s", c.name))
		}
	}

	// sets the endpoint to aws config
	awsconfig := e2e.GetValue("awsconfig")
	if awsconfig == nil {
		awsconfig = aws.NewConfig()
	}
	awsconfig.(*aws.Config).
		WithEndpoint(fmt.Sprintf("http://localhost:%d", edgePort)).
		WithCredentialsChainVerboseErrors(true).
		WithHTTPClient(&http.Client{Timeout: 10 * time.Second}).
		WithMaxRetries(2).
		WithCredentials(credentials.NewStaticCredentials("foo", "bar", "")).
		WithRegion(c.params.Region)
	e2e.AddValue("awsconfig", awsconfig)

	return nil
}
