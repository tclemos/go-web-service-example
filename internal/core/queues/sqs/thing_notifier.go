package sqs

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/pkg/errors"
	"github.com/tclemos/go-dockertest-example/internal/core/domain/events"
)

type ThingNotifier struct {
	sqs       *sqs.SQS
	queueName string
	queueURL  string
	session   *session.Session
}

func NewThingNotifier(qn string, s *session.Session) *ThingNotifier {
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

func (n *ThingNotifier) NotifyThingCreated(e events.ThingCreated) error {

	eventBytes, err := json.Marshal(e)
	thingCreatedEvent := string(eventBytes)

	_, err = n.sqs.SendMessage(&sqs.SendMessageInput{
		DelaySeconds:      aws.Int64(10),
		MessageAttributes: map[string]*sqs.MessageAttributeValue{},
		MessageBody:       &thingCreatedEvent,
		QueueUrl:          &n.queueURL,
	})
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to send message to queue: %s", n.queueName))
	}

	return nil
}
