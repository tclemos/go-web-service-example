package main

import (
	"os"
	"testing"

	"github.com/tclemos/go-dockertest-example/config"
	"github.com/tclemos/go-dockertest-example/e2e"
	"github.com/tclemos/go-dockertest-example/e2e/postgres"
	"github.com/tclemos/go-dockertest-example/e2e/sqs"
)

func TestMain(m *testing.M) {

	c := config.Config{
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

	pc := postgres.NewContainer(c.MyPostgresDb.Database, postgres.Params{
		Port:                c.MyPostgresDb.Port,
		User:                c.MyPostgresDb.User,
		Password:            c.MyPostgresDb.Password,
		Database:            c.MyPostgresDb.Database,
		MigrationsDirectory: "./migrations",
	})

	sc := sqs.NewContainer("aws_sqs", sqs.Params{
		Region: c.ThingNotifier.Region,
		Queues: []sqs.Queue{
			{Name: c.ThingNotifier.QueueName,
				DefaultVisibilityTimeoutInSeconds: 1,
				DelayInSeconds:                    2,
				ReceiveMessageWaitInSeconds:       3,
			},
		},
	})

	e2e.Start(pc, sc)
	e2e.AddValue("config", c)

	code := m.Run()

	e2e.Stop()

	os.Exit(code)
}
