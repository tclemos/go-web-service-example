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

	conn, err := postgres.NewConn(ctx, cfg.MyPostgresDb)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to postgres, err: %w, cfg: %v", err, cfg.MyPostgresDb))
	}
	defer conn.Close(ctx)

	awsconfig := aws.NewConfig()
	session, err := sqs.NewSession(awsconfig)
	if err != nil {
		panic(fmt.Sprintf("Failed to create sqs session, err: %w, cfg: %v", err, awsconfig))
	}

	tr := postgres.NewThingRepository(conn)
	tn := sqs.NewThingNotifier(cfg.ThingNotifier.QueueName, session)
	ts := services.NewThingService(tr, tn)
	tc := http.NewThingsController(ts)

	server := http.NewServer()
	server.AddController(tc)

	server.Start()
}
