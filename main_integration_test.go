package main

import (
	"os"
	"testing"

	"github.com/tclemos/go-dockertest-example/config"
	"github.com/tclemos/go-dockertest-example/e2e"
	"github.com/tclemos/go-dockertest-example/e2e/postgres"
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
	}

	pc := postgres.NewContainer(e2e.Config.MyPostgresDb.Database, postgres.Params{
		Port:                e2e.Config.MyPostgresDb.Port,
		User:                e2e.Config.MyPostgresDb.User,
		Password:            e2e.Config.MyPostgresDb.Password,
		Database:            e2e.Config.MyPostgresDb.Database,
		MigrationsDirectory: "./migrations",
	})

	e2e.StartContainers(pc)

	defer e2e.StopContainers()

	code := m.Run()

	e2e.StopContainers()

	os.Exit(code)
}
