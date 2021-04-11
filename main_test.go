package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/sethvargo/go-retry"
	thingshttpclient "github.com/tclemos/go-web-service-example/adapters/http/client"
	"github.com/tclemos/goit"
	"github.com/tclemos/goit/aws"
	"github.com/tclemos/goit/dockerfile"
	"github.com/tclemos/goit/postgres"
)

const (
	postgresHost     = "localhost"
	postgresPort     = 5432
	postgresUser     = "postgres"
	postgresPassword = "password"
	postgresDatabase = "my_postgres_db"

	awsHost      = "http://localhost"
	awsPort      = 4566
	awsID        = "foo"
	awsSecret    = "bar"
	awsToken     = ""
	awsQueueName = "thing_created"
	awsRegion    = "eu-central-1"

	httpServerHost = "localhost"
	httpServerPort = 8123
)

var sqsService *aws.SqsService

func TestMain(m *testing.M) {

	ctx := context.Background()

	env := map[string]string{
		"THING_APP_POSTGRES_HOST":     postgresHost,
		"THING_APP_POSTGRES_PORT":     strconv.Itoa(postgresPort),
		"THING_APP_POSTGRES_USER":     postgresUser,
		"THING_APP_POSTGRES_PASSWORD": postgresPassword,
		"THING_APP_POSTGRES_DATABASE": postgresDatabase,

		"THING_APP_NOTIFIER_HOST":      awsHost,
		"THING_APP_NOTIFIER_PORT":      strconv.Itoa(awsPort),
		"THING_APP_NOTIFIER_ID":        awsID,
		"THING_APP_NOTIFIER_SECRET":    awsSecret,
		"THING_APP_NOTIFIER_TOKEN":     awsToken,
		"THING_APP_NOTIFIER_QUEUENAME": awsQueueName,
		"THING_APP_NOTIFIER_REGION":    awsRegion,

		"THING_APP_HTTP_SERVER_HOST": httpServerHost,
		"THING_APP_HTTP_SERVER_PORT": strconv.Itoa(httpServerPort),
	}

	postgresContainer := postgres.NewContainer(postgres.Params{
		Port:     postgresPort,
		User:     postgresUser,
		Password: postgresPassword,
		Database: postgresDatabase,
	})

	awsContainer := aws.NewContainer(aws.Params{
		Region: awsRegion,
		Port:   awsPort,
		SqsQueues: []aws.SqsQueue{
			{Name: awsQueueName},
		},
	})

	portBindings := make(map[docker.Port][]docker.PortBinding, 1)
	portBindings[docker.Port(fmt.Sprintf("%d/tcp", httpServerPort))] = []docker.PortBinding{{
		HostIP:   "0.0.0.0",
		HostPort: strconv.Itoa(httpServerPort),
	}}

	appContainer := dockerfile.NewContainer(dockerfile.Params{
		ContainerName:  "thing service",
		DockerFilePath: "./Dockerfile",
		Env:            env,
		PortBindings:   portBindings,
		AfterStart: func(ctx context.Context, r *dockertest.Resource, m *map[string]interface{}) error {
			b, _ := retry.NewFibonacci(500 * time.Millisecond)
			b = retry.WithMaxRetries(10, b)
			b = retry.WithCappedDuration(20*time.Second, b)

			addr := fmt.Sprintf("%s:%d", httpServerHost, httpServerPort)
			c, err := thingshttpclient.NewClient(addr, nil)
			if err != nil {
				return err
			}

			err = retry.Do(ctx, b, func(ctx context.Context) error {
				res, err := c.Ping(ctx)
				if err != nil || res.StatusCode != http.StatusOK {
					fmt.Println("waiting on thing http server to initialize...")
					return retry.RetryableError(err)
				}
				return nil
			})
			if err != nil {
				return err
			}

			return nil
		},
	})

	opt := goit.Options{
		AutoRemoveContainers:         false,
		ExpireContainersAfterSeconds: 600,
	}
	goit.StartWithOptions(ctx, opt, postgresContainer, awsContainer, appContainer)

	sqsService = awsContainer.SqsService

	code := goit.Run(m)

	goit.Stop()

	os.Exit(code)
}

func TestCreateGetThing(t *testing.T) {

	ctx := context.Background()

	addr := fmt.Sprintf("%s:%d", httpServerHost, httpServerPort)
	c, err := thingshttpclient.NewClient(addr, nil)
	if err != nil {
		t.Errorf("failed to create things http client, err: %v", err)
	}

	const code, name = "thingcode", "thingname"

	res, err := c.CreateThing(ctx, thingshttpclient.CreateThingJSONRequestBody{
		Code: code,
		Name: name,
	})
	if err != nil {
		t.Errorf("failed to call create thing API, err: %v", err)
	}

	createdThingRes, err := thingshttpclient.ParseCreateThingResponse(res)
	if err != nil {
		t.Errorf("failed to parse created thing response, err: %v", err)
		return
	}
	if createdThingRes.StatusCode() != http.StatusCreated {
		t.Errorf("invalid status code when creating a thing, expected: %d, found: %d", http.StatusCreated, createdThingRes.StatusCode())
		return
	}

	res, err = c.GetThing(ctx, thingshttpclient.Code(code))
	if err != nil {
		t.Errorf("failed to call get thing API, err: %v", err)
		return
	}
	getThingsRes, err := thingshttpclient.ParseGetThingResponse(res)
	if err != nil {
		t.Errorf("failed to parse created thing response, err: %v", err)
		return
	}
	if getThingsRes.StatusCode() != http.StatusOK {
		t.Errorf("invalid status code when getting a thing, expected: %d, found: %d", http.StatusOK, getThingsRes.StatusCode())
		return
	}
	var thing thingshttpclient.Thing
	err = json.Unmarshal(getThingsRes.Body, &thing)
	if err != nil {
		t.Errorf("failed to unmarshal get thing response body, err: %v", err)
		return
	}

	if thing.Code != code {
		t.Errorf("thing code received is different from the one sent on creation, expected: %s, found: %s", code, thing.Code)
		return
	}

	if thing.Name != name {
		t.Errorf("thing name received is different from the one sent on creation, expected: %s, found: %s", name, thing.Name)
		return
	}
}
