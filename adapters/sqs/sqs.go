package sqs

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/tclemos/go-web-service-example/config"
)

func NewSession(c config.SqsConfig) (*session.Session, error) {
	cfg := aws.NewConfig().
		WithEndpoint(fmt.Sprintf("%s:%d", c.Host, c.Port)).
		WithCredentialsChainVerboseErrors(true).
		WithMaxRetries(2).
		WithCredentials(credentials.NewStaticCredentials(c.Id, c.Secret, c.Token)).
		WithRegion(c.Region).
		WithDisableSSL(true)
	return session.Must(session.NewSession(cfg)), nil
}
