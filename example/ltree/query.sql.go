// Code generated by pggen. DO NOT EDIT.

package ltree

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
	FindTopScienceChildren(ctx context.Context) ([]pgtype.Text, error)

	FindTopScienceChildrenAgg(ctx context.Context) ([]pgtype.Text, error)

	InsertSampleData(ctx context.Context) (pgconn.CommandTag, error)

	FindLtreeInput(ctx context.Context, inLtree pgtype.Text, inLtreeArray []string) (FindLtreeInputRow, error)
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





const findTopScienceChildrenSQL = `SELECT path
FROM test
WHERE path <@ 'Top.Science';`

// FindTopScienceChildren implements Querier.FindTopScienceChildren.
func (q *DBQuerier) FindTopScienceChildren(ctx context.Context) ([]pgtype.Text, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "FindTopScienceChildren")
	rows, err := q.conn.Query(ctx, findTopScienceChildrenSQL)
	if err != nil {
		return nil, fmt.Errorf("query FindTopScienceChildren: %w", err)
	}

	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (pgtype.Text, error) {
  var item pgtype.Text
		if err := row.Scan(&item,
			); err != nil {
			return item, fmt.Errorf("failed to scan: %w", err)
		}
		return item, nil
	})
}

const findTopScienceChildrenAggSQL = `SELECT array_agg(path)
FROM test
WHERE path <@ 'Top.Science';`

// FindTopScienceChildrenAgg implements Querier.FindTopScienceChildrenAgg.
func (q *DBQuerier) FindTopScienceChildrenAgg(ctx context.Context) ([]pgtype.Text, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "FindTopScienceChildrenAgg")
	rows, err := q.conn.Query(ctx, findTopScienceChildrenAggSQL)
	if err != nil {
		return nil, fmt.Errorf("query FindTopScienceChildrenAgg: %w", err)
	}

	return pgx.CollectExactlyOneRow(rows, func(row pgx.CollectableRow) ([]pgtype.Text, error) {
  var item []pgtype.Text
		if err := row.Scan(&item,
			); err != nil {
			return item, fmt.Errorf("failed to scan: %w", err)
		}
		return item, nil
	})
}

const insertSampleDataSQL = `INSERT INTO test
VALUES ('Top'),
       ('Top.Science'),
       ('Top.Science.Astronomy'),
       ('Top.Science.Astronomy.Astrophysics'),
       ('Top.Science.Astronomy.Cosmology'),
       ('Top.Hobbies'),
       ('Top.Hobbies.Amateurs_Astronomy'),
       ('Top.Collections'),
       ('Top.Collections.Pictures'),
       ('Top.Collections.Pictures.Astronomy'),
       ('Top.Collections.Pictures.Astronomy.Stars'),
       ('Top.Collections.Pictures.Astronomy.Galaxies'),
       ('Top.Collections.Pictures.Astronomy.Astronauts');`

// InsertSampleData implements Querier.InsertSampleData.
func (q *DBQuerier) InsertSampleData(ctx context.Context) (pgconn.CommandTag, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "InsertSampleData")
	cmdTag, err := q.conn.Exec(ctx, insertSampleDataSQL)
	if err != nil {
		return pgconn.CommandTag{}, fmt.Errorf("exec query InsertSampleData: %w", err)
	}
	return cmdTag, err
}

const findLtreeInputSQL = `SELECT
  $1::ltree                   AS ltree,
  -- This won't work, but I'm not quite sure why.
  -- Postgres errors with "wrong element type (SQLSTATE 42804)"
  -- All caps because we use regex to find pggen.arg and it confuses pggen.
  -- PGGEN.arg('in_ltree_array_direct')::ltree[]    AS direct_arr,

  -- The parenthesis around the text[] cast are important. They signal to pggen
  -- that we need a text array that Postgres then converts to ltree[].
  ($2::text[])::ltree[] AS text_arr;`

type FindLtreeInputRow struct {
	Ltree   pgtype.Text   `json:"ltree"`
	TextArr []pgtype.Text `json:"text_arr"`
}

// FindLtreeInput implements Querier.FindLtreeInput.
func (q *DBQuerier) FindLtreeInput(ctx context.Context, inLtree pgtype.Text, inLtreeArray []string) (FindLtreeInputRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "FindLtreeInput")
	rows, err := q.conn.Query(ctx, findLtreeInputSQL, inLtree, inLtreeArray)
	if err != nil {
		return FindLtreeInputRow{}, fmt.Errorf("query FindLtreeInput: %w", err)
	}

	return pgx.CollectExactlyOneRow(rows, func(row pgx.CollectableRow) (FindLtreeInputRow, error) {
  var item FindLtreeInputRow
		if err := row.Scan(&item.Ltree, // 'ltree', 'Ltree', 'pgtype.Text', 'github.com/jackc/pgx/v5/pgtype', 'Text'
			&item.TextArr, // 'text_arr', 'TextArr', '[]pgtype.Text', 'github.com/jackc/pgx/v5/pgtype', '[]Text'
			); err != nil {
			return item, fmt.Errorf("failed to scan: %w", err)
		}
		return item, nil
	})
}
