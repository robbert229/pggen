// Code generated by pggen. DO NOT EDIT.

package domain

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

var _ genericConn = (*pgx.Conn)(nil)
var _ RegisterConn = (*pgx.Conn)(nil)

// Querier is a typesafe Go interface backed by SQL queries.
type Querier interface {
	DomainOne(ctx context.Context) (string, error)
}

var _ Querier = &DBQuerier{}

type DBQuerier struct {
	conn  genericConn
}

// genericConn is a connection like *pgx.Conn, pgx.Tx, or *pgxpool.Pool.
type genericConn interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
}

// NewQuerier creates a DBQuerier that implements Querier.
func NewQuerier(ctx context.Context, conn genericConn) (*DBQuerier, error) {
	return &DBQuerier{
		conn: conn, 
	}, nil
}

type typeHook func(ctx context.Context, conn RegisterConn) error

var typeHooks []typeHook

func addHook(hook typeHook) {
	typeHooks = append(typeHooks, hook)
}

type RegisterConn interface {
	LoadType(ctx context.Context, typeName string) (*pgtype.Type, error)
	TypeMap() *pgtype.Map
}

func Register(ctx context.Context, conn RegisterConn) error {
  

	for _, hook := range typeHooks {
		if err := hook(ctx, conn); err != nil {
			return err
		}
	}
	
	return nil
}



const domainOneSQL = `SELECT '90210'::us_postal_code;`

// DomainOne implements Querier.DomainOne.
func (q *DBQuerier) DomainOne(ctx context.Context) (string, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "DomainOne")
	rows, err := q.conn.Query(ctx, domainOneSQL)
	if err != nil {
		return "", fmt.Errorf("query DomainOne: %w", err)
	}

	return pgx.CollectExactlyOneRow(rows, func(row pgx.CollectableRow) (string, error) {
  var item string
		if err := row.Scan(&item,
			); err != nil {
			return item, fmt.Errorf("failed to scan: %w", err)
		}
		return item, nil
	})
}
