package sqs

import "github.com/aws/aws-sdk-go/aws/session"

func NewSession() (*session.Session, error) {
	return session.Must(session.NewSession()), nil
}
