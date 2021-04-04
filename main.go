package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/tclemos/go-web-service-example/actors/http"
	"github.com/tclemos/go-web-service-example/actors/postgres"
	"github.com/tclemos/go-web-service-example/actors/sqs"
	"github.com/tclemos/go-web-service-example/config"
	"github.com/tclemos/go-web-service-example/core/services"
)

func main() {
	// focus on integrated tests
	cfg := config.Config{}
	ctx := context.Background()
	Start(ctx, cfg)
}

func Start(ctx context.Context, cfg config.Config) {

	querier, err := postgres.NewQuerier(ctx, cfg.MyPostgresDb)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to postgres, err: %v, cfg: %v", err, cfg.MyPostgresDb))
	}

	awsconfig := aws.NewConfig()
	session, err := sqs.NewSession(awsconfig)
	if err != nil {
		panic(fmt.Sprintf("Failed to create sqs session, err: %v, cfg: %v", err, awsconfig))
	}

	tr := postgres.NewThingRepository(querier)
	tn := sqs.NewThingNotifier(cfg.ThingNotifier.QueueName, session)
	ts := services.NewThingService(tr, tn)
	tc := http.NewThingsController(ts)

	server := http.NewServer(tc)

	server.Start()
}
