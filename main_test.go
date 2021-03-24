package main

import (
	"context"
	"net/url"
	"os"
	"testing"

	"github.com/tclemos/go-web-service-example/config"
	"github.com/tclemos/goit"
	"github.com/tclemos/goit/aws"
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
		QueueName: "thing_created",
		Region:    "eu-central-1",
	},
}

var dbUrl url.URL
var sqsService *aws.SqsService

func TestMain(m *testing.M) {

	postgresContainer := postgres.NewContainer(postgres.Params{
		Port:     cfg.MyPostgresDb.Port,
		User:     cfg.MyPostgresDb.User,
		Password: cfg.MyPostgresDb.Password,
		Database: cfg.MyPostgresDb.Database,
	})

	awsContainer := aws.NewContainer(aws.Params{
		Region: cfg.ThingNotifier.Region,
		SqsQueues: []aws.SqsQueue{
			{Name: cfg.ThingNotifier.QueueName},
		},
	})

	goit.Start(context.Background(), postgresContainer, awsContainer)
	dbUrl = postgresContainer.Url()
	sqsService = awsContainer.SqsService

	code := m.Run()

	goit.Stop()

	os.Exit(code)
}
