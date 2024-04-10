// Code generated by pggen. DO NOT EDIT.

package function

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"sync"
)

// Querier is a typesafe Go interface backed by SQL queries.
type Querier interface {
	OutParams(ctx context.Context) ([]OutParamsRow, error)
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
func NewQuerier(conn *pgx.Conn) *DBQuerier {
	_ = conn

	return &DBQuerier{
		conn: conn, 
	}
}

// ListItem represents the Postgres composite type "list_item".
type ListItem struct {
	Name  *string `json:"name"`
	Color *string `json:"color"`
}

// ListStats represents the Postgres composite type "list_stats".
type ListStats struct {
	Val1 *string  `json:"val1"`
	Val2 []*int32 `json:"val2"`
}


func register(conn *pgx.Conn){
	//
}



/*type compositeField struct {
	name       string                 // name of the field
	typeName   string                 // Postgres type name
	defaultCodec pgtype.Codec // default value to use
}

func (tr *typeResolver) newCompositeValue(name string, fields ...compositeField) pgtype.Codec {
	if _, codec, ok := tr.findCodec(name); ok {
		return codec
	}

	codecs := make([]pgtype.CompositeCodecField, len(fields))
	isBinaryOk := true
	
	for i, field := range fields {
		oid, codec, ok := tr.findCodec(field.typeName)
		if !ok {
			oid = pgtype.UnknownOID
			codec = field.defaultCodec
		}
		isBinaryOk = isBinaryOk && oid != pgtype.UnknownOID
		
		codecs[i] = pgtype.CompositeCodecField{
			Name: field.name,
			Type: &pgtype.Type{Codec: codec, Name: field.typeName, OID: oid},
		}
	}
	// Okay to ignore error because it's only thrown when the number of field
	// names does not equal the number of ValueTranscoders.
	codec := pgtype.CompositeCodec{Fields: codecs}
	// typ, _ := pgtype.NewCompositeTypeValues(name, fs, codecs)
	// if !isBinaryOk {
	// 	return textPreferrer{ValueTranscoder: typ, typeName: name}
	// }
	return codec
}

func (tr *typeResolver) newArrayValue(name, elemName string, defaultVal func() pgtype.ValueTranscoder) pgtype.Codec {
	if _, val, ok := tr.findCodec(name); ok {
		return val
	}
	
	pgType, ok := tr.pgMap.TypeForName(elemName)
	if !ok {
		panic("unhandled")
	}
	
	return &pgtype.ArrayCodec{ElementType: pgType}
}*/

// newListItem creates a new pgtype.ValueTranscoder for the Postgres
// composite type 'list_item'.
func registernewListItem() pgtype.Codec {
	return tr.newCompositeValue(
		"list_item",
		compositeField{name: "name", typeName: "text", defaultCodec: &pgtype.TextCodec{}},
		compositeField{name: "color", typeName: "text", defaultCodec: &pgtype.TextCodec{}},
	)
}

// newListStats creates a new pgtype.ValueTranscoder for the Postgres
// composite type 'list_stats'.
func registernewListStats() pgtype.Codec {
	return tr.newCompositeValue(
		"list_stats",
		compositeField{name: "val1", typeName: "text", defaultCodec: &pgtype.TextCodec{}},
		compositeField{name: "val2", typeName: "_int4", defaultCodec: &pgtype.Int4Array{}},
	)
}

// newListItemArray creates a new pgtype.Codec for the Postgres
// '_list_item' array type.
func registernewListItemArray() pgtype.Codec {
	return tr.newArrayValue("_list_item", "list_item", tr.newListItem)
}

const outParamsSQL = `SELECT * FROM out_params();`

type OutParamsRow struct {
	Items []ListItem `json:"_items"`
	Stats ListStats  `json:"_stats"`
}

// OutParams implements Querier.OutParams.
func (q *DBQuerier) OutParams(ctx context.Context) ([]OutParamsRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "OutParams")
	rows, err := q.conn.Query(ctx, outParamsSQL)
	if err != nil {
		return nil, fmt.Errorf("query OutParams: %w", err)
	}

	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (OutParamsRow, error) {
		var item OutParamsRow
		if err := row.Scan(
			&item.Items, // '_items', 'Items', '[]ListItem', 'github.com/robbert229/pggen/example/function', '[]ListItem'
			&item.Stats, // '_stats', 'Stats', 'ListStats', 'github.com/robbert229/pggen/example/function', 'ListStats'
		); err != nil {
			return item, fmt.Errorf("failed to scan: %w", err)
		}
		return item, nil
	})
}

type scanCacheKey struct {
	oid      uint32
	format   int16
	typeName string
}

var (
	plans   = make(map[scanCacheKey]pgtype.ScanPlan, 16)
	plansMu sync.RWMutex
)

func planScan(codec pgtype.Codec, fd pgconn.FieldDescription, target any) pgtype.ScanPlan {
	key := scanCacheKey{fd.DataTypeOID, fd.Format, fmt.Sprintf("%T", target)}
	plansMu.RLock()
	plan := plans[key]
	plansMu.RUnlock()
	if plan != nil {
		return plan
	}
	plan = codec.PlanScan(nil, fd.DataTypeOID, fd.Format, target)
	plansMu.Lock()
	plans[key] = plan
	plansMu.Unlock()
	return plan
}

type ptrScanner[T any] struct {
	basePlan pgtype.ScanPlan
}

func (s ptrScanner[T]) Scan(src []byte, dst any) error {
	if src == nil {
		return nil
	}
	d := dst.(**T)
	*d = new(T)
	return s.basePlan.Scan(src, *d)
}

func planPtrScan[T any](codec pgtype.Codec, fd pgconn.FieldDescription, target *T) pgtype.ScanPlan {
	key := scanCacheKey{fd.DataTypeOID, fd.Format, fmt.Sprintf("*%T", target)}
	plansMu.RLock()
	plan := plans[key]
	plansMu.RUnlock()
	if plan != nil {
		return plan
	}
	basePlan := planScan(codec, fd, target)
	ptrPlan := ptrScanner[T]{basePlan}
	plansMu.Lock()
	plans[key] = plan
	plansMu.Unlock()
	return ptrPlan
}