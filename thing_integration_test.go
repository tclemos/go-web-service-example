package main

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/gofrs/uuid"
	"github.com/tclemos/go-dockertest-example/config"
	"github.com/tclemos/go-dockertest-example/e2e"
	"github.com/tclemos/go-dockertest-example/internal/core/queues/sqs"
	"github.com/tclemos/go-dockertest-example/internal/core/services"
	"github.com/tclemos/go-dockertest-example/internal/http"
	"github.com/tclemos/go-dockertest-example/internal/http/requests"
	"github.com/tclemos/go-dockertest-example/internal/repositories/postgres"
)

func TestCreateAndGet(t *testing.T) {

	c := e2e.GetValue("config").(config.Config)
	awsconfig := e2e.GetValue("awsconfig").(*aws.Config)

	conn, err := postgres.NewConn(e2e.Ctx, c.MyPostgresDb)
	if err != nil {
		t.Errorf("Failed to connect to postgres: %v : %w", c.MyPostgresDb, err)
		return
	}
	defer conn.Close(e2e.Ctx)
	tr := postgres.NewThingRepository(conn)

	session, err := sqs.NewSession(c.ThingNotifier.Region)
	if err != nil {
		t.Errorf("Failed to create sqs session: %v : %w", c.ThingNotifier, err)
		return
	}
	tn := sqs.NewThingNotifier(c.ThingNotifier.QueueName, session, awsconfig)

	ts := services.NewThingService(tr, tn)
	tc := http.NewThingsController(ts)

	id, _ := uuid.NewV4()
	code := id.String()
	name := "name"

	createReq := requests.CreateThing{
		Code: code,
		Name: name,
	}
	tc.Create(e2e.Ctx, createReq)

	getReq := requests.GetThing{
		Code: code,
	}
	getRes := tc.Get(e2e.Ctx, getReq)
	if getRes == nil {
		t.Errorf("Thing not found: %s", code)
	} else {
		if getRes.Code != code {
			t.Errorf("Wrong Code. expected: %s. found: %s.", code, getRes.Code)
		}

		if getRes.Name != name {
			t.Errorf("Wrong Name. expected: %s. found: %s.", name, getRes.Name)
		}
	}
}
