package http

import (
	"context"

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

	t := domain.Thing{
		Code: domain.ThingCode(r.Code),
		Name: r.Name,
	}

	err := c.svc.Create(ctx, t)
	if err != nil {
		logger.Errorf(err, "failed to notify thing created")
	}
}

func (c ThingsController) Get(ctx context.Context, r requests.GetThing) *responses.Thing {

	code := domain.ThingCode(r.Code)

	t, err := c.svc.Get(ctx, code)
	if err != nil {
		return nil
	}

	return &responses.Thing{
		Code: t.Code.String(),
		Name: t.Name,
	}
}
