package main

import (
	"testing"

	"github.com/gofrs/uuid"
	"github.com/tclemos/go-dockertest-example/e2e"
	"github.com/tclemos/go-dockertest-example/internal/core/queues/sqs"
	"github.com/tclemos/go-dockertest-example/internal/core/services"
	"github.com/tclemos/go-dockertest-example/internal/http"
	"github.com/tclemos/go-dockertest-example/internal/http/requests"
	"github.com/tclemos/go-dockertest-example/internal/repositories/postgres"
)

func TestCreateAndGet(t *testing.T) {

	conn, err := postgres.NewConn(e2e.Ctx, e2e.Config.MyPostgresDb)
	if err != nil {
		t.Errorf("Failed to connect to postgres: %v : %w", e2e.Config.MyPostgresDb, err)
		return
	}
	defer conn.Close(e2e.Ctx)
	tr := postgres.NewThingRepository(conn)

	session, err := sqs.NewSession()
	if err != nil {
		t.Errorf("Failed to create sqs session: %v : %w", e2e.Config.ThingNotifier, err)
		return
	}
	tn := sqs.NewThingNotifier(e2e.Config.ThingNotifier.QueueName, session)

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
