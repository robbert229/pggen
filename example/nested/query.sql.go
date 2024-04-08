// Code generated by pggen. DO NOT EDIT.

package nested

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
	ArrayNested2(ctx context.Context) ([]ProductImageType, error)

	Nested3(ctx context.Context) ([]ProductImageSetType, error)
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

// Dimensions represents the Postgres composite type "dimensions".
type Dimensions struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

// ProductImageSetType represents the Postgres composite type "product_image_set_type".
type ProductImageSetType struct {
	Name      string             `json:"name"`
	OrigImage ProductImageType   `json:"orig_image"`
	Images    []ProductImageType `json:"images"`
}

// ProductImageType represents the Postgres composite type "product_image_type".
type ProductImageType struct {
	Source     string     `json:"source"`
	Dimensions Dimensions `json:"dimensions"`
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

// newDimensions creates a new pgtype.ValueTranscoder for the Postgres
// composite type 'dimensions'.
func registernewDimensions() pgtype.Codec {
	return tr.newCompositeValue(
		"dimensions",
		compositeField{name: "width", typeName: "int4", defaultCodec: &pgtype.Int4Codec{}},
		compositeField{name: "height", typeName: "int4", defaultCodec: &pgtype.Int4Codec{}},
	)
}

// newProductImageSetType creates a new pgtype.ValueTranscoder for the Postgres
// composite type 'product_image_set_type'.
func registernewProductImageSetType() pgtype.Codec {
	return tr.newCompositeValue(
		"product_image_set_type",
		compositeField{name: "name", typeName: "text", defaultCodec: &pgtype.TextCodec{}},
		compositeField{name: "orig_image", typeName: "product_image_type", defaultCodec: tr.newProductImageType()},
		compositeField{name: "images", typeName: "_product_image_type", defaultCodec: tr.newProductImageTypeArray()},
	)
}

// newProductImageType creates a new pgtype.ValueTranscoder for the Postgres
// composite type 'product_image_type'.
func registernewProductImageType() pgtype.Codec {
	return tr.newCompositeValue(
		"product_image_type",
		compositeField{name: "source", typeName: "text", defaultCodec: &pgtype.TextCodec{}},
		compositeField{name: "dimensions", typeName: "dimensions", defaultCodec: tr.newDimensions()},
	)
}

// newProductImageTypeArray creates a new pgtype.Codec for the Postgres
// '_product_image_type' array type.
func registernewProductImageTypeArray() pgtype.Codec {
	return tr.newArrayValue("_product_image_type", "product_image_type", tr.newProductImageType)
}

const arrayNested2SQL = `SELECT
  ARRAY [
    ROW ('img2', ROW (22, 22)::dimensions)::product_image_type,
    ROW ('img3', ROW (33, 33)::dimensions)::product_image_type
    ] AS images;`

// ArrayNested2 implements Querier.ArrayNested2.
func (q *DBQuerier) ArrayNested2(ctx context.Context) ([]ProductImageType, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "ArrayNested2")
	rows, err := q.conn.Query(ctx, arrayNested2SQL)
	if err != nil {
		return nil, fmt.Errorf("query ArrayNested2: %w", err)
	}

	return pgx.CollectExactlyOneRow(rows, func(row pgx.CollectableRow) ([]ProductImageType, error) {
		var item []ProductImageType
		if err := row.Scan(
			&item,
		); err != nil {
			return item, fmt.Errorf("failed to scan: %w", err)
		}
		return item, nil
	})
}

const nested3SQL = `SELECT
  ROW (
    'name', -- name
    ROW ('img1', ROW (11, 11)::dimensions)::product_image_type, -- orig_image
    ARRAY [ --images
      ROW ('img2', ROW (22, 22)::dimensions)::product_image_type,
      ROW ('img3', ROW (33, 33)::dimensions)::product_image_type
      ]
    )::product_image_set_type;`

// Nested3 implements Querier.Nested3.
func (q *DBQuerier) Nested3(ctx context.Context) ([]ProductImageSetType, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "Nested3")
	rows, err := q.conn.Query(ctx, nested3SQL)
	if err != nil {
		return nil, fmt.Errorf("query Nested3: %w", err)
	}

	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (ProductImageSetType, error) {
		var item ProductImageSetType
		if err := row.Scan(
			&item,
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