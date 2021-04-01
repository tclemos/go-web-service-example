package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/tclemos/go-web-service-example/actors/postgres/db"
	"github.com/tclemos/go-web-service-example/core/domain"
	"github.com/tclemos/go-web-service-example/core/port"
)

// ThingRepository for a postgres db
type ThingRepository struct {
	q db.Querier
}

// NewThingRepository creates and initialises an instance of
func NewThingRepository(q db.Querier) port.ThingRepository {
	return &ThingRepository{
		q: q,
	}
}

// Create a thing
func (r ThingRepository) Create(ctx context.Context, t domain.Thing) error {
	arg := db.CreateThingParams{
		ID:   int32(t.ID),
		Code: uuid.New(),
		Name: t.Name,
	}
	_, err := r.q.CreateThing(ctx, arg)
	return err
}

// Get a thing
func (r ThingRepository) Get(ctx context.Context, id domain.ThingID) (*domain.Thing, error) {
	t, err := r.q.GetThing(ctx, int32(id))
	if err != nil {
		return nil, err
	}

	return &domain.Thing{
		ID:   domain.ThingID(t.ID),
		Code: domain.ThingCode(t.Code),
		Name: t.Name,
	}, nil
}

// GetByCode a thing
func (r ThingRepository) GetByCode(ctx context.Context, code domain.ThingCode) (*domain.Thing, error) {

	t, err := r.q.GetThingByCode(ctx, uuid.UUID(code))
	if err != nil {
		return nil, err
	}

	return &domain.Thing{
		ID:   domain.ThingID(t.ID),
		Code: domain.ThingCode(t.Code),
		Name: t.Name,
	}, nil
}

// Update a thing
func (r ThingRepository) Update(ctx context.Context, t domain.Thing) error {

	arg := db.UpdateThingParams{
		ID:   int32(t.ID),
		Code: uuid.UUID(t.Code),
		Name: t.Name,
	}
	_, err := r.q.UpdateThing(ctx, arg)
	return err
}

// Delete a thing
func (r ThingRepository) Delete(ctx context.Context, id domain.ThingID) error {
	err := r.q.DeleteThing(ctx, int32(id))
	return err
}

// DeleteByCode a thing
func (r ThingRepository) DeleteByCode(ctx context.Context, code domain.ThingCode) error {
	err := r.q.DeleteThingByCode(ctx, uuid.UUID(code))
	return err
}
