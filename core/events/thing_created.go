package events

import "github.com/tclemos/go-web-service-example/core/domain"

type ThingCreated struct {
	Thing domain.Thing `json:"thing"`
}
