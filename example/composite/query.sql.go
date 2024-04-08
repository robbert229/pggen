// Code generated by pggen. DO NOT EDIT.

package composite

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
	SearchScreenshots(ctx context.Context, params SearchScreenshotsParams) ([]SearchScreenshotsRow, error)

	SearchScreenshotsOneCol(ctx context.Context, params SearchScreenshotsOneColParams) ([][]Blocks, error)

	InsertScreenshotBlocks(ctx context.Context, screenshotID int, body string) (InsertScreenshotBlocksRow, error)

	ArraysInput(ctx context.Context, arrays Arrays) (Arrays, error)

	UserEmails(ctx context.Context) (UserEmail, error)
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

// Arrays represents the Postgres composite type "arrays".
type Arrays struct {
	Texts  []string   `json:"texts"`
	Int8s  []*int     `json:"int8s"`
	Bools  []bool     `json:"bools"`
	Floats []*float64 `json:"floats"`
}

// Blocks represents the Postgres composite type "blocks".
type Blocks struct {
	ID           int    `json:"id"`
	ScreenshotID int    `json:"screenshot_id"`
	Body         string `json:"body"`
}

// UserEmail represents the Postgres composite type "user_email".
type UserEmail struct {
	ID    string      `json:"id"`
	Email pgtype.Text `json:"email"`
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

// newArrays creates a new pgtype.ValueTranscoder for the Postgres
// composite type 'arrays'.
func registernewArrays() pgtype.Codec {
	return tr.newCompositeValue(
		"arrays",
		compositeField{name: "texts", typeName: "_text", defaultCodec: &pgtype.TextArray{}},
		compositeField{name: "int8s", typeName: "_int8", defaultCodec: &pgtype.Int8Array{}},
		compositeField{name: "bools", typeName: "_bool", defaultCodec: &pgtype.BoolArray{}},
		compositeField{name: "floats", typeName: "_float8", defaultCodec: &pgtype.Float8Array{}},
	)
}

// newArraysInit creates an initialized pgtype.ValueTranscoder for the
// Postgres composite type 'arrays' to encode query parameters.
func registernewArraysInit(v Arrays) pgtype.Codec {
	return tr.setCodec(tr.newArrays(), tr.newArraysRaw(v))
}

// newArraysRaw returns all composite fields for the Postgres composite
// type 'arrays' as a slice of interface{} to encode query parameters.
func registernewArraysRaw(v Arrays) []interface{} {
	return []interface{}{
		v.Texts,
		v.Int8s,
		v.Bools,
		v.Floats,
	}
}

// newBlocks creates a new pgtype.ValueTranscoder for the Postgres
// composite type 'blocks'.
func registernewBlocks() pgtype.Codec {
	return tr.newCompositeValue(
		"blocks",
		compositeField{name: "id", typeName: "int4", defaultCodec: &pgtype.Int4Codec{}},
		compositeField{name: "screenshot_id", typeName: "int8", defaultCodec: &pgtype.Int8Codec{}},
		compositeField{name: "body", typeName: "text", defaultCodec: &pgtype.TextCodec{}},
	)
}

// newUserEmail creates a new pgtype.ValueTranscoder for the Postgres
// composite type 'user_email'.
func registernewUserEmail() pgtype.Codec {
	return tr.newCompositeValue(
		"user_email",
		compositeField{name: "id", typeName: "text", defaultCodec: &pgtype.TextCodec{}},
		compositeField{name: "email", typeName: "citext", defaultCodec: &pgtype.TextCodec{}},
	)
}

// newBlocksArray creates a new pgtype.Codec for the Postgres
// '_blocks' array type.
func registernewBlocksArray() pgtype.Codec {
	return tr.newArrayValue("_blocks", "blocks", tr.newBlocks)
}

// newboolArrayRaw returns all elements for the Postgres array type '_bool'
// as a slice of interface{} for use with the pgtype.Value Set method.
func registernewboolArrayRaw(vs []bool) []interface{} {
	elems := make([]interface{}, len(vs))
	for i, v := range vs {
		elems[i] = v
	}
	return elems
}

const searchScreenshotsSQL = `SELECT
  ss.id,
  array_agg(bl) AS blocks
FROM screenshots ss
  JOIN blocks bl ON bl.screenshot_id = ss.id
WHERE bl.body LIKE $1 || '%'
GROUP BY ss.id
ORDER BY ss.id
LIMIT $2 OFFSET $3;`

type SearchScreenshotsParams struct {
	Body   string `json:"Body"`
	Limit  int    `json:"Limit"`
	Offset int    `json:"Offset"`
}

type SearchScreenshotsRow struct {
	ID     int      `json:"id"`
	Blocks []Blocks `json:"blocks"`
}

// SearchScreenshots implements Querier.SearchScreenshots.
func (q *DBQuerier) SearchScreenshots(ctx context.Context, params SearchScreenshotsParams) ([]SearchScreenshotsRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "SearchScreenshots")
	rows, err := q.conn.Query(ctx, searchScreenshotsSQL, params.Body, params.Limit, params.Offset)
	if err != nil {
		return nil, fmt.Errorf("query SearchScreenshots: %w", err)
	}

	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (SearchScreenshotsRow, error) {
		var item SearchScreenshotsRow
		if err := row.Scan(
			&item.ID, // 'id', 'ID', 'int', '', 'int'
			&item.Blocks, // 'blocks', 'Blocks', '[]Blocks', 'github.com/robbert229/pggen/example/composite', '[]Blocks'
		); err != nil {
			return item, fmt.Errorf("failed to scan: %w", err)
		}
		return item, nil
	})
}

const searchScreenshotsOneColSQL = `SELECT
  array_agg(bl) AS blocks
FROM screenshots ss
  JOIN blocks bl ON bl.screenshot_id = ss.id
WHERE bl.body LIKE $1 || '%'
GROUP BY ss.id
ORDER BY ss.id
LIMIT $2 OFFSET $3;`

type SearchScreenshotsOneColParams struct {
	Body   string `json:"Body"`
	Limit  int    `json:"Limit"`
	Offset int    `json:"Offset"`
}

// SearchScreenshotsOneCol implements Querier.SearchScreenshotsOneCol.
func (q *DBQuerier) SearchScreenshotsOneCol(ctx context.Context, params SearchScreenshotsOneColParams) ([][]Blocks, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "SearchScreenshotsOneCol")
	rows, err := q.conn.Query(ctx, searchScreenshotsOneColSQL, params.Body, params.Limit, params.Offset)
	if err != nil {
		return nil, fmt.Errorf("query SearchScreenshotsOneCol: %w", err)
	}

	return pgx.CollectRows(rows, func(row pgx.CollectableRow) ([]Blocks, error) {
		var item []Blocks
		if err := row.Scan(
			&item,
		); err != nil {
			return item, fmt.Errorf("failed to scan: %w", err)
		}
		return item, nil
	})
}

const insertScreenshotBlocksSQL = `WITH screens AS (
  INSERT INTO screenshots (id) VALUES ($1)
    ON CONFLICT DO NOTHING
)
INSERT
INTO blocks (screenshot_id, body)
VALUES ($1, $2)
RETURNING id, screenshot_id, body;`

type InsertScreenshotBlocksRow struct {
	ID           int    `json:"id"`
	ScreenshotID int    `json:"screenshot_id"`
	Body         string `json:"body"`
}

// InsertScreenshotBlocks implements Querier.InsertScreenshotBlocks.
func (q *DBQuerier) InsertScreenshotBlocks(ctx context.Context, screenshotID int, body string) (InsertScreenshotBlocksRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "InsertScreenshotBlocks")
	rows, err := q.conn.Query(ctx, insertScreenshotBlocksSQL, screenshotID, body)
	if err != nil {
		return InsertScreenshotBlocksRow{}, fmt.Errorf("query InsertScreenshotBlocks: %w", err)
	}

	return pgx.CollectExactlyOneRow(rows, func(row pgx.CollectableRow) (InsertScreenshotBlocksRow, error) {
		var item InsertScreenshotBlocksRow
		if err := row.Scan(
			&item.ID, // 'id', 'ID', 'int', '', 'int'
			&item.ScreenshotID, // 'screenshot_id', 'ScreenshotID', 'int', '', 'int'
			&item.Body, // 'body', 'Body', 'string', '', 'string'
		); err != nil {
			return item, fmt.Errorf("failed to scan: %w", err)
		}
		return item, nil
	})
}

const arraysInputSQL = `SELECT $1::arrays;`

// ArraysInput implements Querier.ArraysInput.
func (q *DBQuerier) ArraysInput(ctx context.Context, arrays Arrays) (Arrays, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "ArraysInput")
	rows, err := q.conn.Query(ctx, arraysInputSQL, q.types.newArraysInit(arrays))
	if err != nil {
		return Arrays{}, fmt.Errorf("query ArraysInput: %w", err)
	}

	return pgx.CollectExactlyOneRow(rows, func(row pgx.CollectableRow) (Arrays, error) {
		var item Arrays
		if err := row.Scan(
			&item,
		); err != nil {
			return item, fmt.Errorf("failed to scan: %w", err)
		}
		return item, nil
	})
}

const userEmailsSQL = `SELECT ('foo', 'bar@example.com')::user_email;`

// UserEmails implements Querier.UserEmails.
func (q *DBQuerier) UserEmails(ctx context.Context) (UserEmail, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "UserEmails")
	rows, err := q.conn.Query(ctx, userEmailsSQL)
	if err != nil {
		return UserEmail{}, fmt.Errorf("query UserEmails: %w", err)
	}

	return pgx.CollectExactlyOneRow(rows, func(row pgx.CollectableRow) (UserEmail, error) {
		var item UserEmail
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