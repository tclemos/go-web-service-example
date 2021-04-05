package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Config struct {
	MyPostgresDb  PostgresConfig `json:"my_postgres_db"`
	ThingNotifier SqsConfig      `json:"thing_notifier"`
	Http          HttpConfig     `json:"http"`
}

type PostgresConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
	Schema   string `json:"schema"`
}

type SqsConfig struct {
	Region    string `json:"region"`
	Host      string `json:"host"`
	Port      int    `json:"port"`
	Id        string `json:"id"`
	Secret    string `json:"secret"`
	Token     string `json:"token"`
	QueueName string `json:"queue_name"`
}

type HttpConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
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
