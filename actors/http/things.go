package http

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
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

func (c ThingsController) SetupHandlers(s ServerInterface) {

}

// (GET /things)
func (c ThingsController) FindThing(ctx echo.Context, params FindThingParams) error {
	return nil
}

// (POST /things)
func (c ThingsController) CreateThing(ctx echo.Context) error {
	var t CreateThingJSONBody
	err := ctx.Bind(&t)
	if err != nil {
		return httpError(ctx, http.StatusBadRequest, "Invalid request body")
	}

	cd, err := uuid.Parse(t.Code)
	if err != nil {
		logger.Errorf(err, "invalid thing code")
		return err
	}

	newThing := domain.Thing{
		Code: domain.ThingCode(cd),
		Name: t.Name,
	}

	if err := c.svc.Create(ctx.Request().Context(), newThing); err != nil {
		logger.Errorf(err, "failed to notify thing created")
	}

	return ctx.String(http.StatusCreated, "")
}

// (PUT /things)
func (c ThingsController) UpdateThing(ctx echo.Context) error {
	return nil
}

// (DELETE /things/{code})
func (c ThingsController) DeleteThing(ctx echo.Context, code Code) error {
	return nil
}

// (GET /things/{code})
func (c ThingsController) GetThingsCode(ctx echo.Context, code Code) error {
	cd, err := uuid.Parse(string(code))
	if err != nil {
		return httpError(ctx, http.StatusBadRequest, fmt.Sprintf("invalid thing code: %s", code))
	}

	t, err := c.svc.GetByCode(ctx.Request().Context(), domain.ThingCode(cd))
	if err != nil {
		return nil
	}

	res := &Thing{
		Code: t.Code.String(),
		Name: t.Name,
	}

	return ctx.JSON(http.StatusOK, res)
}
