package sqs

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

func NewSession() (*session.Session, error) {

	host := os.Getenv("THING_APP_NOTIFIER_HOST")
	port := os.Getenv("THING_APP_NOTIFIER_PORT")
	id := os.Getenv("THING_APP_NOTIFIER_ID")
	secret := os.Getenv("THING_APP_NOTIFIER_SECRET")
	token := os.Getenv("THING_APP_NOTIFIER_TOKEN")
	region := os.Getenv("THING_APP_NOTIFIER_REGION")

	cfg := aws.NewConfig().
		WithEndpoint(fmt.Sprintf("%s:%d", host, port)).
		WithCredentialsChainVerboseErrors(true).
		WithMaxRetries(2).
		WithCredentials(credentials.NewStaticCredentials(id, secret, token)).
		WithRegion(region).
		WithDisableSSL(true)
	return session.Must(session.NewSession(cfg)), nil
}
