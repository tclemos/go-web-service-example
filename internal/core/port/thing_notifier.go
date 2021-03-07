package port

import (
	"github.com/tclemos/go-dockertest-example/internal/core/events"
)

type ThingNotifier interface {
	NotifyThingCreated(events.ThingCreated) error
}
