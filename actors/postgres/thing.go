package postgres

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/tclemos/go-web-service-example/core/domain"
	"github.com/tclemos/go-web-service-example/core/port"
)

// ThingRepository for a postgres db
type ThingRepository struct {
	conn *pgx.Conn
}

// NewThingRepository creates and initialises an instance of
func NewThingRepository(c *pgx.Conn) port.ThingRepository {
	return &ThingRepository{
		conn: c,
	}
}

// Create a thing
func (r ThingRepository) Create(ctx context.Context, t domain.Thing) error {
	_, err := r.conn.Exec(ctx, "insert into things(code, name) values($1,$2)", t.Code.String(), t.Name)
	return err
}

// Get a thing
func (r ThingRepository) Get(ctx context.Context, code domain.ThingCode) (*domain.Thing, error) {
	var c, n string
	err := r.conn.
		QueryRow(ctx, "select code, name from things where code=$1", code.String()).
		Scan(&c, &n)
	if err != nil {
		return nil, err
	}

	return &domain.Thing{
		Code: domain.ThingCode(c),
		Name: n,
	}, nil
}

// Update a thing
func (r ThingRepository) Update(ctx context.Context, t domain.Thing) error {
	_, err := r.conn.Exec(ctx, "update things set name=$2 where code=$1", t.Code.String(), t.Name)
	return err
}

// Delete a thing
func (r ThingRepository) Delete(ctx context.Context, code domain.ThingCode) error {
	_, err := r.conn.Exec(context.Background(), "delete from things where code=$1", code.String())
	return err
}
