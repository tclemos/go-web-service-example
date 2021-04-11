package main

import (
	"context"
	"fmt"

	"github.com/spf13/viper"
	"github.com/tclemos/go-web-service-example/adapters/http"
	"github.com/tclemos/go-web-service-example/adapters/http/controllers"
	"github.com/tclemos/go-web-service-example/adapters/postgres"
	"github.com/tclemos/go-web-service-example/adapters/sqs"
	"github.com/tclemos/go-web-service-example/core/services"
)

func main() {
	ctx := context.Background()
	viper.SetEnvPrefix("THING_APP_")

	querier, err := postgres.NewQuerier(ctx)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to postgres, err: %v", err))
	}

	tr := postgres.NewThingRepository(querier)
	tn := sqs.NewThingNotifier()
	ts := services.NewThingService(tr, tn)
	tc := controllers.NewThingsController(ts)

	server := http.NewServer(tc)

	server.Start()
}
