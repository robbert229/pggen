// Code generated by pggen. DO NOT EDIT.

package void

import (
	"context"
	"database/sql/driver"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

var _ genericConn = (*pgx.Conn)(nil)
var _ RegisterConn = (*pgx.Conn)(nil)

// Querier is a typesafe Go interface backed by SQL queries.
type Querier interface {
	VoidOnly(ctx context.Context) (pgconn.CommandTag, error)

	VoidOnlyTwoParams(ctx context.Context, id int32) (pgconn.CommandTag, error)

	VoidTwo(ctx context.Context) (string, error)

	VoidThree(ctx context.Context) (VoidThreeRow, error)

	VoidThree2(ctx context.Context) ([]string, error)
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
  
    conn.TypeMap().RegisterType(&pgtype.Type{
      Name: "void",
      OID: 2278,
      Codec: voidCodec{},
    })
  

	for _, hook := range typeHooks {
		if err := hook(ctx, conn); err != nil {
			return err
		}
	}
	
	return nil
}


type voidCodec struct {}

var _ pgtype.Codec = &voidCodec{}

// FormatSupported returns true if the format is supported.
func (voidCodec)  FormatSupported(int16) bool {
	return true
}

// PreferredFormat returns the preferred format.
func (voidCodec) PreferredFormat() int16 {
	return pgtype.TextFormatCode
}

// PlanEncode returns an EncodePlan for encoding value into PostgreSQL format for oid and format. If no plan can be
// found then nil is returned.
func (voidCodec) PlanEncode(m *pgtype.Map, oid uint32, format int16, value any) pgtype.EncodePlan {
	return nil
}

// PlanScan returns a ScanPlan for scanning a PostgreSQL value into a destination with the same type as target. If
// no plan can be found then nil is returned.
func (voidCodec) PlanScan(m *pgtype.Map, oid uint32, format int16, target any) pgtype.ScanPlan {
	return nil
}

// DecodeDatabaseSQLValue returns src decoded into a value compatible with the sql.Scanner interface.
func (voidCodec) DecodeDatabaseSQLValue(m *pgtype.Map, oid uint32, format int16, src []byte) (driver.Value, error) {
	return nil, nil
}

// DecodeValue returns src decoded into its default format.
func (voidCodec) DecodeValue(m *pgtype.Map, oid uint32, format int16, src []byte) (any, error) {
	return nil, nil
}

type pgVoid struct {}

// Scan implements the database/sql Scanner interface.
func (dst *pgVoid) Scan(src any) error {
	return nil
}

// Value implements the database/sql/driver Valuer interface.
func (src pgVoid) Value() (driver.Value, error) {
	return nil, nil
}


const voidOnlySQL = `SELECT void_fn();`

// VoidOnly implements Querier.VoidOnly.
func (q *DBQuerier) VoidOnly(ctx context.Context) (pgconn.CommandTag, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "VoidOnly")
	cmdTag, err := q.conn.Exec(ctx, voidOnlySQL)
	if err != nil {
		return pgconn.CommandTag{}, fmt.Errorf("exec query VoidOnly: %w", err)
	}
	return cmdTag, err
}

const voidOnlyTwoParamsSQL = `SELECT void_fn_two_params($1, 'text');`

// VoidOnlyTwoParams implements Querier.VoidOnlyTwoParams.
func (q *DBQuerier) VoidOnlyTwoParams(ctx context.Context, id int32) (pgconn.CommandTag, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "VoidOnlyTwoParams")
	cmdTag, err := q.conn.Exec(ctx, voidOnlyTwoParamsSQL, id)
	if err != nil {
		return pgconn.CommandTag{}, fmt.Errorf("exec query VoidOnlyTwoParams: %w", err)
	}
	return cmdTag, err
}

const voidTwoSQL = `SELECT void_fn(), 'foo' as name;`

// VoidTwo implements Querier.VoidTwo.
func (q *DBQuerier) VoidTwo(ctx context.Context) (string, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "VoidTwo")
	rows, err := q.conn.Query(ctx, voidTwoSQL)
	if err != nil {
		return "", fmt.Errorf("query VoidTwo: %w", err)
	}

	return pgx.CollectExactlyOneRow(rows, func(row pgx.CollectableRow) (string, error) {
  var item string
		if err := row.Scan(
			&pgVoid{},
			&item,
			); err != nil {
			return item, fmt.Errorf("failed to scan: %w", err)
		}
		return item, nil
	})
}

const voidThreeSQL = `SELECT void_fn(), 'foo' as foo, 'bar' as bar;`

type VoidThreeRow struct {
	Foo string `json:"foo"`
	Bar string `json:"bar"`
}

// VoidThree implements Querier.VoidThree.
func (q *DBQuerier) VoidThree(ctx context.Context) (VoidThreeRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "VoidThree")
	rows, err := q.conn.Query(ctx, voidThreeSQL)
	if err != nil {
		return VoidThreeRow{}, fmt.Errorf("query VoidThree: %w", err)
	}

	return pgx.CollectExactlyOneRow(rows, func(row pgx.CollectableRow) (VoidThreeRow, error) {
  var item VoidThreeRow
		if err := row.Scan(
			&pgVoid{},
			&item.Foo, // 'foo', 'Foo', 'string', '', 'string'
			&item.Bar, // 'bar', 'Bar', 'string', '', 'string'
			); err != nil {
			return item, fmt.Errorf("failed to scan: %w", err)
		}
		return item, nil
	})
}

const voidThree2SQL = `SELECT 'foo' as foo, void_fn(), void_fn();`

// VoidThree2 implements Querier.VoidThree2.
func (q *DBQuerier) VoidThree2(ctx context.Context) ([]string, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "VoidThree2")
	rows, err := q.conn.Query(ctx, voidThree2SQL)
	if err != nil {
		return nil, fmt.Errorf("query VoidThree2: %w", err)
	}

	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (string, error) {
  var item string
		if err := row.Scan(&item,
			
			&pgVoid{},
			
			&pgVoid{},
			); err != nil {
			return item, fmt.Errorf("failed to scan: %w", err)
		}
		return item, nil
	})
}
