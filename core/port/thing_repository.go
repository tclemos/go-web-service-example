package port

import (
	"context"

	"github.com/tclemos/go-web-service-example/core/domain"
)

type ThingRepository interface {
	Create(context.Context, domain.Thing) error
	Get(context.Context, domain.ThingID) (*domain.Thing, error)
	GetByCode(context.Context, domain.ThingCode) (*domain.Thing, error)
	Update(context.Context, domain.Thing) error
	DeleteByCode(context.Context, domain.ThingCode) error
}
