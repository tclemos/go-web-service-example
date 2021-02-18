package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Config struct {
	MyPostgresDb  PostgresConfig `json:"my_postgres_db"`
	ThingNotifier SqsConfig      `json:"thing_notifier"`
}

type PostgresConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}

type SqsConfig struct {
	QueueName string `json:"queue_name`
}

func LoadConfig(path string) Config {
	c := Config{}
	j, err := ioutil.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("%s: %w", "failed to load configuration", err))
	}
	json.Unmarshal(j, &c)
	return c
}
