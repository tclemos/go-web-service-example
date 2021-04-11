// Code generated by sqlc. DO NOT EDIT.

package db

import (
	"context"
	"database/sql"
	"fmt"
)

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func New(db DBTX) *Queries {
	return &Queries{db: db}
}

func Prepare(ctx context.Context, db DBTX) (*Queries, error) {
	q := Queries{db: db}
	var err error
	if q.createAnotherThingStmt, err = db.PrepareContext(ctx, createAnotherThing); err != nil {
		return nil, fmt.Errorf("error preparing query CreateAnotherThing: %w", err)
	}
	if q.createThingStmt, err = db.PrepareContext(ctx, createThing); err != nil {
		return nil, fmt.Errorf("error preparing query CreateThing: %w", err)
	}
	if q.deleteAnotherThingStmt, err = db.PrepareContext(ctx, deleteAnotherThing); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteAnotherThing: %w", err)
	}
	if q.deleteAnotherThingByCodeStmt, err = db.PrepareContext(ctx, deleteAnotherThingByCode); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteAnotherThingByCode: %w", err)
	}
	if q.deleteThingStmt, err = db.PrepareContext(ctx, deleteThing); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteThing: %w", err)
	}
	if q.deleteThingByCodeStmt, err = db.PrepareContext(ctx, deleteThingByCode); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteThingByCode: %w", err)
	}
	if q.getAnotherThingStmt, err = db.PrepareContext(ctx, getAnotherThing); err != nil {
		return nil, fmt.Errorf("error preparing query GetAnotherThing: %w", err)
	}
	if q.getAnotherThingByCodeStmt, err = db.PrepareContext(ctx, getAnotherThingByCode); err != nil {
		return nil, fmt.Errorf("error preparing query GetAnotherThingByCode: %w", err)
	}
	if q.getThingStmt, err = db.PrepareContext(ctx, getThing); err != nil {
		return nil, fmt.Errorf("error preparing query GetThing: %w", err)
	}
	if q.getThingByCodeStmt, err = db.PrepareContext(ctx, getThingByCode); err != nil {
		return nil, fmt.Errorf("error preparing query GetThingByCode: %w", err)
	}
	if q.listAnotherThingsStmt, err = db.PrepareContext(ctx, listAnotherThings); err != nil {
		return nil, fmt.Errorf("error preparing query ListAnotherThings: %w", err)
	}
	if q.listThingsStmt, err = db.PrepareContext(ctx, listThings); err != nil {
		return nil, fmt.Errorf("error preparing query ListThings: %w", err)
	}
	if q.updateAnotherThingStmt, err = db.PrepareContext(ctx, updateAnotherThing); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateAnotherThing: %w", err)
	}
	if q.updateThingStmt, err = db.PrepareContext(ctx, updateThing); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateThing: %w", err)
	}
	return &q, nil
}

func (q *Queries) Close() error {
	var err error
	if q.createAnotherThingStmt != nil {
		if cerr := q.createAnotherThingStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createAnotherThingStmt: %w", cerr)
		}
	}
	if q.createThingStmt != nil {
		if cerr := q.createThingStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createThingStmt: %w", cerr)
		}
	}
	if q.deleteAnotherThingStmt != nil {
		if cerr := q.deleteAnotherThingStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteAnotherThingStmt: %w", cerr)
		}
	}
	if q.deleteAnotherThingByCodeStmt != nil {
		if cerr := q.deleteAnotherThingByCodeStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteAnotherThingByCodeStmt: %w", cerr)
		}
	}
	if q.deleteThingStmt != nil {
		if cerr := q.deleteThingStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteThingStmt: %w", cerr)
		}
	}
	if q.deleteThingByCodeStmt != nil {
		if cerr := q.deleteThingByCodeStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteThingByCodeStmt: %w", cerr)
		}
	}
	if q.getAnotherThingStmt != nil {
		if cerr := q.getAnotherThingStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getAnotherThingStmt: %w", cerr)
		}
	}
	if q.getAnotherThingByCodeStmt != nil {
		if cerr := q.getAnotherThingByCodeStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getAnotherThingByCodeStmt: %w", cerr)
		}
	}
	if q.getThingStmt != nil {
		if cerr := q.getThingStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getThingStmt: %w", cerr)
		}
	}
	if q.getThingByCodeStmt != nil {
		if cerr := q.getThingByCodeStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getThingByCodeStmt: %w", cerr)
		}
	}
	if q.listAnotherThingsStmt != nil {
		if cerr := q.listAnotherThingsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing listAnotherThingsStmt: %w", cerr)
		}
	}
	if q.listThingsStmt != nil {
		if cerr := q.listThingsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing listThingsStmt: %w", cerr)
		}
	}
	if q.updateAnotherThingStmt != nil {
		if cerr := q.updateAnotherThingStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateAnotherThingStmt: %w", cerr)
		}
	}
	if q.updateThingStmt != nil {
		if cerr := q.updateThingStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateThingStmt: %w", cerr)
		}
	}
	return err
}

func (q *Queries) exec(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (sql.Result, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).ExecContext(ctx, args...)
	case stmt != nil:
		return stmt.ExecContext(ctx, args...)
	default:
		return q.db.ExecContext(ctx, query, args...)
	}
}

func (q *Queries) query(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (*sql.Rows, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryContext(ctx, args...)
	default:
		return q.db.QueryContext(ctx, query, args...)
	}
}

func (q *Queries) queryRow(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) *sql.Row {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryRowContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryRowContext(ctx, args...)
	default:
		return q.db.QueryRowContext(ctx, query, args...)
	}
}

type Queries struct {
	db                           DBTX
	tx                           *sql.Tx
	createAnotherThingStmt       *sql.Stmt
	createThingStmt              *sql.Stmt
	deleteAnotherThingStmt       *sql.Stmt
	deleteAnotherThingByCodeStmt *sql.Stmt
	deleteThingStmt              *sql.Stmt
	deleteThingByCodeStmt        *sql.Stmt
	getAnotherThingStmt          *sql.Stmt
	getAnotherThingByCodeStmt    *sql.Stmt
	getThingStmt                 *sql.Stmt
	getThingByCodeStmt           *sql.Stmt
	listAnotherThingsStmt        *sql.Stmt
	listThingsStmt               *sql.Stmt
	updateAnotherThingStmt       *sql.Stmt
	updateThingStmt              *sql.Stmt
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db:                           tx,
		tx:                           tx,
		createAnotherThingStmt:       q.createAnotherThingStmt,
		createThingStmt:              q.createThingStmt,
		deleteAnotherThingStmt:       q.deleteAnotherThingStmt,
		deleteAnotherThingByCodeStmt: q.deleteAnotherThingByCodeStmt,
		deleteThingStmt:              q.deleteThingStmt,
		deleteThingByCodeStmt:        q.deleteThingByCodeStmt,
		getAnotherThingStmt:          q.getAnotherThingStmt,
		getAnotherThingByCodeStmt:    q.getAnotherThingByCodeStmt,
		getThingStmt:                 q.getThingStmt,
		getThingByCodeStmt:           q.getThingByCodeStmt,
		listAnotherThingsStmt:        q.listAnotherThingsStmt,
		listThingsStmt:               q.listThingsStmt,
		updateAnotherThingStmt:       q.updateAnotherThingStmt,
		updateThingStmt:              q.updateThingStmt,
	}
}
