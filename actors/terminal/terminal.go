package main

import (
	"context"
	"os"

	"github.com/google/uuid"
	"github.com/tclemos/go-web-service-example/core/domain"
	"github.com/tclemos/go-web-service-example/core/services"
)

func main() {

	// thing -f "./things"

	var svc services.ThingService

	if os.Args[1] == "--create" {
		code, err := uuid.Parse(os.Args[2])
		if err != nil {
			panic(err)
		}
		name := os.Args[3]

		t := domain.Thing{
			Code: domain.ThingCode(code),
			Name: name,
		}

		svc.Create(context.Background(), t)
	}

}
