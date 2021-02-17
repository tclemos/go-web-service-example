package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Config struct {
	MyPostgresDb PostgresConfig `json:"my_postgres_db"`
}

type PostgresConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
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