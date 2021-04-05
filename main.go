package main

import (
	"context"
	"fmt"

	"github.com/tclemos/go-web-service-example/actors/http"
	"github.com/tclemos/go-web-service-example/actors/http/controllers"
	"github.com/tclemos/go-web-service-example/actors/postgres"
	"github.com/tclemos/go-web-service-example/actors/sqs"
	"github.com/tclemos/go-web-service-example/config"
	"github.com/tclemos/go-web-service-example/core/services"
)

func main() {
	// focus on integrated tests
	ctx := context.Background()

	cfg := config.LoadConfig("./config/config.json")

	Start(ctx, cfg)
}

func Start(ctx context.Context, cfg config.Config) {

	querier, err := postgres.NewQuerier(ctx, cfg.MyPostgresDb)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to postgres, err: %v, cfg: %v", err, cfg.MyPostgresDb))
	}

	tr := postgres.NewThingRepository(querier)
	tn := sqs.NewThingNotifier(cfg.ThingNotifier.QueueName, cfg.ThingNotifier)
	ts := services.NewThingService(tr, tn)
	tc := controllers.NewThingsController(ts)

	server := http.NewServer(tc)

	server.Start()
}
