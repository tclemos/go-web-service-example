package port

import (
	"github.com/tclemos/go-web-service-example/core/events"
)

type ThingNotifier interface {
	NotifyThingCreated(events.ThingCreated) error
}
