package http

import (
	"context"

	"github.com/google/uuid"
	"github.com/tclemos/go-web-service-example/actors/http/requests"
	"github.com/tclemos/go-web-service-example/actors/http/responses"
	"github.com/tclemos/go-web-service-example/actors/logger"
	"github.com/tclemos/go-web-service-example/core/domain"
	"github.com/tclemos/go-web-service-example/core/services"
)

type ThingsController struct {
	svc *services.ThingService
}

func NewThingsController(ts *services.ThingService) *ThingsController {
	return &ThingsController{
		svc: ts,
	}
}

func (c ThingsController) Create(ctx context.Context, r requests.CreateThing) {
	cd, err := uuid.Parse(r.Code)
	if err != nil {
		logger.Errorf(err, "invalid thing code")
		return
	}

	t := domain.Thing{
		Code: domain.ThingCode(cd),
		Name: r.Name,
	}

	if err := c.svc.Create(ctx, t); err != nil {
		logger.Errorf(err, "failed to notify thing created")
	}
}

func (c ThingsController) Get(ctx context.Context, r requests.GetThing) *responses.Thing {

	cd, err := uuid.Parse(r.Code)
	if err != nil {
		logger.Errorf(err, "invalid thing code")
		return nil
	}

	t, err := c.svc.GetByCode(ctx, domain.ThingCode(cd))
	if err != nil {
		return nil
	}

	return &responses.Thing{
		Code: t.Code.String(),
		Name: t.Name,
	}
}
