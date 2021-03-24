package main

import (
	"context"
	"os"

	"github.com/tclemos/go-web-service-example/core/domain"
	"github.com/tclemos/go-web-service-example/core/services"
)

func main() {

	// thing -f "./things"

	var svc services.ThingService

	if os.Args[1] == "--create" {
		code := os.Args[2]
		name := os.Args[3]

		t := domain.Thing{
			Code: domain.ThingCode(code),
			Name: name,
		}

		svc.Create(context.Background(), t)
	}

}
