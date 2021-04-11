package ports

import "github.com/tclemos/go-web-service-example/core/domain"

type ThingNotifier interface {
	NotifyThingCreated(domain.ThingCreated) error
}
