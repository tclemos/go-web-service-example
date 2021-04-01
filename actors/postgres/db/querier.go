// Code generated by sqlc. DO NOT EDIT.

package db

import (
	"context"

	"github.com/google/uuid"
)

type Querier interface {
	CreateThing(ctx context.Context, arg CreateThingParams) (Thing, error)
	DeleteThing(ctx context.Context, id int32) error
	DeleteThingByCode(ctx context.Context, code uuid.UUID) error
	GetThing(ctx context.Context, id int32) (Thing, error)
	GetThingByCode(ctx context.Context, code uuid.UUID) (Thing, error)
	ListThings(ctx context.Context) ([]Thing, error)
	UpdateThing(ctx context.Context, arg UpdateThingParams) (Thing, error)
}

var _ Querier = (*Queries)(nil)