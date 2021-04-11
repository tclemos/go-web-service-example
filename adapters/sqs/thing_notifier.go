package sqs

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/pkg/errors"
	"github.com/tclemos/go-web-service-example/config"
	"github.com/tclemos/go-web-service-example/core/domain"
)

type ThingNotifier struct {
	sqs       *sqs.SQS
	queueName string
	queueURL  string
	session   *session.Session
}

type thing struct {
	Code   string `json:"code"`
	Name   string `json:"name"`
	Status string `json:"status"`
}
type thingCreatedMessage struct {
	Thing thing `json:"thing"`
}

func NewThingNotifier(qn string, c config.SqsConfig) *ThingNotifier {

	s, err := NewSession(c)
	if err != nil {
		panic(fmt.Sprintf("Failed to create thing notifier, err: %v, cfg: %v", err, c))
	}

	svc := sqs.New(s)
	urlOutput, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: &qn,
	})

	if err != nil {
		panic(fmt.Errorf("failed to get the queueUrl of the queue: %s. err: %v", qn, err))
	}

	return &ThingNotifier{
		sqs:       svc,
		queueName: qn,
		queueURL:  *urlOutput.QueueUrl,
		session:   s,
	}
}

func (n *ThingNotifier) NotifyThingCreated(e domain.ThingCreated) error {

	b, err := json.Marshal(thingCreatedMessage{
		Thing: thing{
			Code:   e.Thing.Code.String(),
			Name:   e.Thing.Name,
			Status: string(e.Thing.Status),
		},
	})
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to marshal thing created: %v", e))
	}

	message := string(b)
	_, err = n.sqs.SendMessage(&sqs.SendMessageInput{
		MessageAttributes: map[string]*sqs.MessageAttributeValue{},
		MessageBody:       &message,
		QueueUrl:          &n.queueURL,
	})
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to send message to queue: %s", n.queueName))
	}

	return nil
}
