package sqs

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/pkg/errors"
)

const (
	queuePort = 9324
	httpPort  = 9325
)

const (
	queueConfigTemplate = `	%s {
		defaultVisibilityTimeout = %d seconds
		delay = %d seconds
		receiveMessageWait = %d seconds
	}`

	configTemplate = `include classpath("application.conf")                                            
                                           
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
}`
)

type Queue struct {
	Name                              string
	DefaultVisibilityTimeoutInSeconds int
	DelayInSeconds                    int
	ReceiveMessageWaitInSeconds       int
}

// Params needed to start a postgres container
type Params struct {
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

	// create config folder
	if _, err := r.Exec([]string{"mkdir", "config"}, dockertest.ExecOptions{}); err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to create config directory to sqs container: %s", c.name))
	}

	// create config file
	config := c.createConfig()
	if _, err := r.Exec([]string{"bash", "-c", fmt.Sprintf("echo '%s' > ./config/elasticmq.conf", config)}, dockertest.ExecOptions{}); err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to write config to sqs container: %s", c.name))
	}

	// restart elastic mq
	if _, err := r.Exec([]string{"bash", "-c", "supervisorctl restart elasticmq"}, dockertest.ExecOptions{}); err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to restart elastic mq on sqs container: %s", c.name))
	}
	return nil
}

func (c *Container) createConfig() string {
	var queuesConfig string
	for _, q := range c.params.Queues {
		if strings.TrimSpace(q.Name) == "" {
			panic("Queue name is required")
		}

		queuesConfig = queuesConfig + fmt.Sprintln(fmt.Sprintf(queueConfigTemplate,
			q.Name, q.DefaultVisibilityTimeoutInSeconds, q.DelayInSeconds, q.ReceiveMessageWaitInSeconds))
	}

	config := fmt.Sprintf(configTemplate,
		httpPort, httpPort, queuesConfig)

	return config
}
