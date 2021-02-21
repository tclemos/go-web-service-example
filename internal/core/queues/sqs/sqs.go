package sqs

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

func NewSession(c *aws.Config) (*session.Session, error) {
	return session.Must(session.NewSession(c)), nil
}
