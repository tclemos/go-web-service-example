package main

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/gofrs/uuid"
	"github.com/tclemos/go-dockertest-example/config"
	"github.com/tclemos/go-dockertest-example/e2e"
	"github.com/tclemos/go-dockertest-example/e2e/localstack"
	"github.com/tclemos/go-dockertest-example/internal/core/domain/events"
	"github.com/tclemos/go-dockertest-example/internal/core/queues/sqs"
	"github.com/tclemos/go-dockertest-example/internal/core/services"
	"github.com/tclemos/go-dockertest-example/internal/http"
	"github.com/tclemos/go-dockertest-example/internal/http/requests"
	"github.com/tclemos/go-dockertest-example/internal/repositories/postgres"
)

func TestCreateAndGet(t *testing.T) {

	// Arrange
	c := e2e.GetValue("config").(config.Config)
	awsconfig := e2e.GetValue("awsconfig").(*aws.Config)

	conn, err := postgres.NewConn(e2e.Ctx, c.MyPostgresDb)
	if err != nil {
		t.Errorf("Failed to connect to postgres: %v : %w", c.MyPostgresDb, err)
		return
	}
	defer conn.Close(e2e.Ctx)
	tr := postgres.NewThingRepository(conn)

	session, err := sqs.NewSession(awsconfig)
	if err != nil {
		t.Errorf("Failed to create sqs session: %v : %w", c.ThingNotifier, err)
		return
	}
	sqsReceiver := localstack.NewSqsReceiver(c.ThingNotifier.QueueName, session)
	tn := sqs.NewThingNotifier(c.ThingNotifier.QueueName, session)

	ts := services.NewThingService(tr, tn)
	tc := http.NewThingsController(ts)

	id, _ := uuid.NewV4()
	code := id.String()
	name := "name"

	// Ack
	createReq := requests.CreateThing{
		Code: code,
		Name: name,
	}
	tc.Create(e2e.Ctx, createReq)

	// Assert Thing Created
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

	// awaits 1 second until the message can be received
	// Assert Thing Created notification
	time.Sleep(10000)
	messages, err := sqsReceiver.Receive()
	if err != nil {
		t.Errorf("failed to receive thing created message, err: %w", err)
	}

	qtyMessages := len(messages)
	if qtyMessages != 1 {
		t.Errorf("wrong number of messages found, expected: 1, found: %d", qtyMessages)
	}
	firsMessage := messages[0]
	thingCreatedJSON := firsMessage.Body
	var thingCreated events.ThingCreated
	err = json.Unmarshal([]byte(*thingCreatedJSON), &thingCreated)
	if err != nil {
		t.Errorf("failed to read thing created message body as json, body: %s, err: %w", *thingCreatedJSON, err)
	}

	if thingCreated.Thing.Code.String() != code {
		t.Errorf("Wrong Thing Created Code. expected: %s. found: %s.", code, thingCreated.Thing.Code.String())
	}
	if thingCreated.Thing.Name != name {
		t.Errorf("Wrong Thing Created Name. expected: %s. found: %s.", name, thingCreated.Thing.Name)
	}
}
