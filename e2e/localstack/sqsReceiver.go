package localstack

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

var (
	visibilityTimeout = int64(60)
)

type SqsReceiver struct {
	sqs       *sqs.SQS
	queueName string
	queueURL  string
	session   *session.Session
}

func NewSqsReceiver(qn string, s *session.Session) *SqsReceiver {
	svc := sqs.New(s)
	urlOutput, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: &qn,
	})

	if err != nil {
		panic(fmt.Errorf("failed to get the queueUrl of the queue: %s. err: %v", qn, err))
	}

	return &SqsReceiver{
		sqs:       svc,
		queueName: qn,
		queueURL:  *urlOutput.QueueUrl,
		session:   s,
	}
}

func (r SqsReceiver) Receive() ([]*sqs.Message, error) {
	msgResult, err := r.sqs.ReceiveMessage(&sqs.ReceiveMessageInput{
		AttributeNames: []*string{
			aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
		},
		MessageAttributeNames: []*string{
			aws.String(sqs.QueueAttributeNameAll),
		},
		QueueUrl:            &r.queueURL,
		MaxNumberOfMessages: aws.Int64(1),
		VisibilityTimeout:   &visibilityTimeout,
	})
	if err != nil {
		return nil, err
	}

	return msgResult.Messages, nil
}
