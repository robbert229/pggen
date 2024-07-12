// Code generated by pggen. DO NOT EDIT.

package go_pointer_types

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
	GenSeries1(ctx context.Context) (*int, error)

	GenSeries(ctx context.Context) ([]*int, error)

	GenSeriesArr1(ctx context.Context) ([]int, error)

	GenSeriesArr(ctx context.Context) ([][]int, error)

	GenSeriesStr1(ctx context.Context) (*string, error)

	GenSeriesStr(ctx context.Context) ([]*string, error)
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



const genSeries1SQL = `SELECT n
FROM generate_series(0, 2) n
LIMIT 1;`

// GenSeries1 implements Querier.GenSeries1.
func (q *DBQuerier) GenSeries1(ctx context.Context) (*int, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "GenSeries1")
	rows, err := q.conn.Query(ctx, genSeries1SQL)
	if err != nil {
		return nil, fmt.Errorf("query GenSeries1: %w", err)
	}

	return pgx.CollectExactlyOneRow(rows, func(row pgx.CollectableRow) (*int, error) {
  var item *int
		if err := row.Scan(&item,
			); err != nil {
			return item, fmt.Errorf("failed to scan: %w", err)
		}
		return item, nil
	})
}

const genSeriesSQL = `SELECT n
FROM generate_series(0, 2) n;`

// GenSeries implements Querier.GenSeries.
func (q *DBQuerier) GenSeries(ctx context.Context) ([]*int, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "GenSeries")
	rows, err := q.conn.Query(ctx, genSeriesSQL)
	if err != nil {
		return nil, fmt.Errorf("query GenSeries: %w", err)
	}

	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (*int, error) {
  var item *int
		if err := row.Scan(&item,
			); err != nil {
			return item, fmt.Errorf("failed to scan: %w", err)
		}
		return item, nil
	})
}

const genSeriesArr1SQL = `SELECT array_agg(n)
FROM generate_series(0, 2) n;`

// GenSeriesArr1 implements Querier.GenSeriesArr1.
func (q *DBQuerier) GenSeriesArr1(ctx context.Context) ([]int, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "GenSeriesArr1")
	rows, err := q.conn.Query(ctx, genSeriesArr1SQL)
	if err != nil {
		return nil, fmt.Errorf("query GenSeriesArr1: %w", err)
	}

	return pgx.CollectExactlyOneRow(rows, func(row pgx.CollectableRow) ([]int, error) {
  var item []int
		if err := row.Scan(&item,
			); err != nil {
			return item, fmt.Errorf("failed to scan: %w", err)
		}
		return item, nil
	})
}

const genSeriesArrSQL = `SELECT array_agg(n)
FROM generate_series(0, 2) n;`

// GenSeriesArr implements Querier.GenSeriesArr.
func (q *DBQuerier) GenSeriesArr(ctx context.Context) ([][]int, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "GenSeriesArr")
	rows, err := q.conn.Query(ctx, genSeriesArrSQL)
	if err != nil {
		return nil, fmt.Errorf("query GenSeriesArr: %w", err)
	}

	return pgx.CollectRows(rows, func(row pgx.CollectableRow) ([]int, error) {
  var item []int
		if err := row.Scan(&item,
			); err != nil {
			return item, fmt.Errorf("failed to scan: %w", err)
		}
		return item, nil
	})
}

const genSeriesStr1SQL = `SELECT n::text
FROM generate_series(0, 2) n
LIMIT 1;`

// GenSeriesStr1 implements Querier.GenSeriesStr1.
func (q *DBQuerier) GenSeriesStr1(ctx context.Context) (*string, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "GenSeriesStr1")
	rows, err := q.conn.Query(ctx, genSeriesStr1SQL)
	if err != nil {
		return nil, fmt.Errorf("query GenSeriesStr1: %w", err)
	}

	return pgx.CollectExactlyOneRow(rows, func(row pgx.CollectableRow) (*string, error) {
  var item *string
		if err := row.Scan(&item,
			); err != nil {
			return item, fmt.Errorf("failed to scan: %w", err)
		}
		return item, nil
	})
}

const genSeriesStrSQL = `SELECT n::text
FROM generate_series(0, 2) n;`

// GenSeriesStr implements Querier.GenSeriesStr.
func (q *DBQuerier) GenSeriesStr(ctx context.Context) ([]*string, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "GenSeriesStr")
	rows, err := q.conn.Query(ctx, genSeriesStrSQL)
	if err != nil {
		return nil, fmt.Errorf("query GenSeriesStr: %w", err)
	}

	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (*string, error) {
  var item *string
		if err := row.Scan(&item,
			); err != nil {
			return item, fmt.Errorf("failed to scan: %w", err)
		}
		return item, nil
	})
}
