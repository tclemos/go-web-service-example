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

	e2e.Config = config.Config{
		MyPostgresDb: config.PostgresConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "postgres",
			Password: "password",
			Database: "my_postgres_db",
		},
		ThingNotifier: config.SqsConfig{
			QueueName: "thing_created",
		},
	}

	pc := postgres.NewContainer(e2e.Config.MyPostgresDb.Database, postgres.Params{
		Port:                e2e.Config.MyPostgresDb.Port,
		User:                e2e.Config.MyPostgresDb.User,
		Password:            e2e.Config.MyPostgresDb.Password,
		Database:            e2e.Config.MyPostgresDb.Database,
		MigrationsDirectory: "./migrations",
	})

	sc := sqs.NewContainer("aws_sqs", sqs.Params{
		Queues: []sqs.Queue{
			{Name: e2e.Config.ThingNotifier.QueueName},
		},
	})

	e2e.Start(pc, sc)

	code := m.Run()

	e2e.Stop()

	os.Exit(code)
}
