package events

import "github.com/tclemos/go-dockertest-example/internal/core/domain"

type ThingCreated struct {
	Thing domain.Thing `json:"thing"`
}
