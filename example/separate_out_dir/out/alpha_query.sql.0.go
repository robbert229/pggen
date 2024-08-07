// Code generated by pggen. DO NOT EDIT.

package out

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
	AlphaNested(ctx context.Context) (string, error)

	AlphaCompositeArray(ctx context.Context) ([]Alpha, error)

	Alpha(ctx context.Context) (string, error)

	Bravo(ctx context.Context) (string, error)
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



// Alpha represents the Postgres composite type "alpha".
type Alpha struct {
	Key *string `json:"key"`
}




	// codec_newAlpha is a codec for the composite type of the same name
	func codec_newAlpha(conn RegisterConn) (pgtype.Codec, error) {
		
		    field0, ok := conn.TypeMap().TypeForName("text")
			if !ok {
				return nil, fmt.Errorf("type not found: text")
			}
		
		
		return &pgtype.CompositeCodec{
			Fields: []pgtype.CompositeCodecField{
				
					{
						Name: "key",
						Type: field0,
					},
				
			},
		}, nil
	}

	func register_newAlpha(
		ctx context.Context,
		conn RegisterConn,
	) error {
		t, err := conn.LoadType(
			ctx,
			"\"alpha\"",
		)
		if err != nil {
			return fmt.Errorf("newAlpha failed to load type: %w", err)
		}
		
		conn.TypeMap().RegisterType(t)

		return nil
	}

	func init(){
		addHook(register_newAlpha) 
	}
	


	// codec_newAlphaArray is a codec for the composite type of the same name
	func codec_newAlphaArray(conn RegisterConn) (pgtype.Codec, error) {
		elementType, ok := conn.TypeMap().TypeForName("alpha")
		if !ok {
			return nil, fmt.Errorf("type not found: alpha")
		}

		return &pgtype.ArrayCodec{
			ElementType: elementType,
		}, nil
	}

	func register_newAlphaArray(
		ctx context.Context,
		conn RegisterConn,
	) error {
		t, err := conn.LoadType(
			ctx,
			"\"_alpha\"",
		)
		if err != nil {
			return fmt.Errorf("newAlphaArray failed to load type: %w", err)
		}

		conn.TypeMap().RegisterType(t)

		return nil
	}

	func init(){
		addHook(register_newAlphaArray) 
	}
	

const alphaNestedSQL = `SELECT 'alpha_nested' as output;`

// AlphaNested implements Querier.AlphaNested.
func (q *DBQuerier) AlphaNested(ctx context.Context) (string, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "AlphaNested")
	rows, err := q.conn.Query(ctx, alphaNestedSQL)
	if err != nil {
		return "", fmt.Errorf("query AlphaNested: %w", err)
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

const alphaCompositeArraySQL = `SELECT ARRAY[ROW('key')]::alpha[];`

// AlphaCompositeArray implements Querier.AlphaCompositeArray.
func (q *DBQuerier) AlphaCompositeArray(ctx context.Context) ([]Alpha, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "AlphaCompositeArray")
	rows, err := q.conn.Query(ctx, alphaCompositeArraySQL)
	if err != nil {
		return nil, fmt.Errorf("query AlphaCompositeArray: %w", err)
	}

	return pgx.CollectExactlyOneRow(rows, func(row pgx.CollectableRow) ([]Alpha, error) {
  var item []Alpha
		if err := row.Scan(&item,
			); err != nil {
			return item, fmt.Errorf("failed to scan: %w", err)
		}
		return item, nil
	})
}
