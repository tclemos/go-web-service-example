package main

import (
	"os"
	"testing"

	"github.com/tclemos/go-dockertest-example/config"
	"github.com/tclemos/go-dockertest-example/e2e"
	"github.com/tclemos/go-dockertest-example/e2e/localstack"
	"github.com/tclemos/go-dockertest-example/e2e/postgres"
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

	myPostgresDb := postgres.NewContainer(c.MyPostgresDb.Database, postgres.Params{
		Port:                c.MyPostgresDb.Port,
		User:                c.MyPostgresDb.User,
		Password:            c.MyPostgresDb.Password,
		Database:            c.MyPostgresDb.Database,
		MigrationsDirectory: "./migrations",
	})

	localstack := localstack.NewContainer("localstack", localstack.Params{
		Region: c.ThingNotifier.Region,
		SqsQueues: []localstack.SqsQueue{
			{Name: c.ThingNotifier.QueueName},
		},
	})

	e2e.Start(myPostgresDb, localstack)
	e2e.AddValue("config", c)

	code := m.Run()

	e2e.Stop()

	os.Exit(code)
}
