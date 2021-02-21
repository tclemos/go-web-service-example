package localstack

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/pkg/errors"
	"github.com/sethvargo/go-retry"
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
		WithRegion(c.params.Region).
		WithDisableSSL(true)

	e2e.AddValue("awsconfig", awsconfig)

	s, err := session.NewSession(awsconfig.(*aws.Config))
	if err != nil {
		fmt.Println("localstack: waiting on server to start...")
		return err
	}
	svc := sqs.New(s)

	// await initialization
	c.awaitInitialization(ctx, svc)

	// create sqs queues
	for _, q := range c.params.SqsQueues {
		_, err := svc.CreateQueue(&sqs.CreateQueueInput{
			QueueName: aws.String(q.Name),
		})
		if err != nil {
			fmt.Printf("localstack: failed to create queue: %s\n", q.Name)
			return err
		}
	}

	return nil
}

func (c *Container) awaitInitialization(ctx context.Context, svc *sqs.SQS) error {
	// prepare a connection verification interval. Use a Fibonacci backoff
	// instead of exponential so wait times scale appropriately.
	b, err := retry.NewFibonacci(500 * time.Millisecond)
	if err != nil {
		err = errors.Wrap(err, "failed to configure retries to wait initialization")
		return err
	}

	b = retry.WithMaxRetries(10, b)
	b = retry.WithCappedDuration(100*time.Second, b)

	// Tries to create a queue to make sure the localstack is up.
	var createQueue *sqs.CreateQueueOutput
	err = retry.Do(ctx, b, func(ctx context.Context) error {
		createQueue, err = svc.CreateQueue(&sqs.CreateQueueInput{
			QueueName: aws.String("test-Resource"),
		})
		if err != nil {
			fmt.Println("localstack: waiting on server to initialize...")
			return retry.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		err = errors.Wrap(err, "failed to start localstack")
		return err
	}

	if _, err := svc.DeleteQueue(&sqs.DeleteQueueInput{
		QueueUrl: createQueue.QueueUrl,
	}); err != nil {
		return err
	}

	return nil
}
