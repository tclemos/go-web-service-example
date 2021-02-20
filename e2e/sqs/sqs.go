package sqs

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/pkg/errors"
	"github.com/tclemos/go-dockertest-example/e2e"
)

const (
	queuePort = 9324
	httpPort  = 9325
)

type Queue struct {
	Name                              string
	DefaultVisibilityTimeoutInSeconds int
	DelayInSeconds                    int
	ReceiveMessageWaitInSeconds       int
}

// Params needed to start a postgres container
type Params struct {
	Region string
	Queues []Queue
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

// Options to start a postgres container accordingly to the params
func (c *Container) Options() (*dockertest.RunOptions, error) {
	pb := map[docker.Port][]docker.PortBinding{}
	pb[docker.Port(fmt.Sprintf("%d/tcp", queuePort))] = []docker.PortBinding{{HostIP: "0.0.0.0", HostPort: strconv.Itoa(queuePort)}}
	pb[docker.Port(fmt.Sprintf("%d/tcp", httpPort))] = []docker.PortBinding{{HostIP: "0.0.0.0", HostPort: strconv.Itoa(httpPort)}}
	return &dockertest.RunOptions{
		Name:         c.name,
		Repository:   "roribio16/alpine-sqs",
		Tag:          "1.2.0",
		Env:          []string{},
		PortBindings: pb,
	}, nil
}

// AfterStart will check the connection and execute migrations
func (c *Container) AfterStart(ctx context.Context, r *dockertest.Resource) error {
	// create config files
	config := c.createConfig()
	if _, err := r.Exec([]string{"bash", "-c", fmt.Sprintf("echo '%s' > /opt/config/elasticmq.conf", config)}, dockertest.ExecOptions{}); err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to write config to sqs container: %s", c.name))
	}
	insight := c.createInsight()
	if _, err := r.Exec([]string{"bash", "-c", fmt.Sprintf("echo '%s' > /opt/config/sqs-insight.conf", insight)}, dockertest.ExecOptions{}); err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to write insight to sqs container: %s", c.name))
	}

	// restart elastic mq
	if _, err := r.Exec([]string{"bash", "-c", "supervisorctl restart elasticmq"}, dockertest.ExecOptions{}); err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to restart elastic mq on sqs container: %s", c.name))
	}

	// sets the endpoint to aws config
	awsconfig := e2e.GetValue("awsconfig")
	if awsconfig == nil {
		awsconfig = aws.NewConfig()
	}
	awsconfig.(*aws.Config).WithEndpoint(fmt.Sprintf("http://localhost:%d", queuePort))
	e2e.AddValue("awsconfig", awsconfig)

	return nil
}

func (c *Container) createConfig() string {
	var queuesConfig string
	for _, q := range c.params.Queues {
		if strings.TrimSpace(q.Name) == "" {
			panic("Queue name is required")
		}

		queuesConfig = queuesConfig + fmt.Sprintln(fmt.Sprintf(
			`	%s {
		defaultVisibilityTimeout = %d seconds
		delay = %d seconds
		receiveMessageWait = %d seconds
	}`,
			q.Name, q.DefaultVisibilityTimeoutInSeconds, q.DelayInSeconds, q.ReceiveMessageWaitInSeconds))
	}

	config := fmt.Sprintf(
		`include classpath("application.conf")

node-address {
	protocol = http
	host = "*"
	port = %d
	context-path = ""
}

rest-sqs {
	enabled = true
	bind-port = %d
	bind-hostname = "0.0.0.0"
	// Possible values: relaxed, strict
	sqs-limits = strict
}

queues {
%s
}`, queuePort, queuePort, queuesConfig)

	return config
}

func (c *Container) createInsight() string {
	config := fmt.Sprintf(
		`{
	"port": %d,
	"rememberMessages": 100,

	"endpoints": [],

	"dynamicEndpoints": [
		{
			"key": "notValidKey",
			"secretKey": "notValidSecret",
			"region": "%s",
			"url": "http://localhost:%d",
			"visibility": 0
		}
	]
}`, httpPort, c.params.Region, queuePort)

	return config
}
