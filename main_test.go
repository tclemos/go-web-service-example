package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/sethvargo/go-retry"
	thingshttpclient "github.com/tclemos/go-web-service-example/adapters/http/client"
	"github.com/tclemos/go-web-service-example/config"
	"github.com/tclemos/goit"
	"github.com/tclemos/goit/aws"
	"github.com/tclemos/goit/dockerfile"
	"github.com/tclemos/goit/postgres"
)

var cfg = config.Config{
	MyPostgresDb: config.PostgresConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "password",
		Database: "my_postgres_db",
	},
	ThingNotifier: config.SqsConfig{
		Host:      "http://localhost",
		Port:      4566,
		Id:        "foo",
		Secret:    "bar",
		Token:     "",
		QueueName: "thing_created",
		Region:    "eu-central-1",
	},
	Http: config.HttpConfig{
		Host: "localhost",
		Port: 8123,
	},
}

var dbUrl url.URL
var sqsService *aws.SqsService

func TestMain(m *testing.M) {

	ctx := context.Background()

	postgresContainer := postgres.NewContainer(postgres.Params{
		Port:     cfg.MyPostgresDb.Port,
		User:     cfg.MyPostgresDb.User,
		Password: cfg.MyPostgresDb.Password,
		Database: cfg.MyPostgresDb.Database,
	})

	awsContainer := aws.NewContainer(aws.Params{
		Region: cfg.ThingNotifier.Region,
		Port:   cfg.ThingNotifier.Port,
		SqsQueues: []aws.SqsQueue{
			{Name: cfg.ThingNotifier.QueueName},
		},
	})

	portBindings := make(map[docker.Port][]docker.PortBinding, 1)
	portBindings[docker.Port(fmt.Sprintf("%d/tcp", cfg.Http.Port))] = []docker.PortBinding{{
		HostIP:   "0.0.0.0",
		HostPort: strconv.Itoa(cfg.Http.Port),
	}}

	appContainer := dockerfile.NewContainer(dockerfile.Params{
		ContainerName:  "thing service",
		DockerFilePath: "./Dockerfile",
		PortBindings:   portBindings,
		AfterStart: func(c context.Context, r *dockertest.Resource, m *map[string]interface{}) error {
			b, _ := retry.NewFibonacci(500 * time.Millisecond)
			b = retry.WithMaxRetries(10, b)
			b = retry.WithCappedDuration(20*time.Second, b)

			err := retry.Do(c, b, func(ctx context.Context) error {
				addr := fmt.Sprintf("http://localhost:%d/ping", cfg.Http.Port)
				res, err := http.DefaultClient.Get(addr)
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

	goit.Start(ctx, postgresContainer, awsContainer, appContainer)
	dbUrl = postgresContainer.Url()
	sqsService = awsContainer.SqsService

	code := goit.Run(m)

	goit.Stop()

	os.Exit(code)
}

func TestCreateGetThing(t *testing.T) {

	ctx := context.Background()

	addr := fmt.Sprintf("%s:%d", cfg.Http.Host, cfg.Http.Port)
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
