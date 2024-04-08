// Code generated by pggen. DO NOT EDIT.

package enums

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
	FindAllDevices(ctx context.Context) ([]FindAllDevicesRow, error)

	InsertDevice(ctx context.Context, mac pgtype.Macaddr, typePg DeviceType) (pgconn.CommandTag, error)

	// Select an array of all device_type enum values.
	FindOneDeviceArray(ctx context.Context) ([]DeviceType, error)

	// Select many rows of device_type enum values.
	FindManyDeviceArray(ctx context.Context) ([][]DeviceType, error)

	// Select many rows of device_type enum values with multiple output columns.
	FindManyDeviceArrayWithNum(ctx context.Context) ([]FindManyDeviceArrayWithNumRow, error)

	// Regression test for https://github.com/robbert229/pggen/issues/23.
	EnumInsideComposite(ctx context.Context) (Device, error)
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

// Device represents the Postgres composite type "device".
type Device struct {
	Mac  pgtype.Macaddr `json:"mac"`
	Type DeviceType     `json:"type"`
}

// newDeviceTypeEnum creates a new pgtype.ValueTranscoder for the
// Postgres enum type 'device_type'.
func registernewDeviceTypeEnum() pgtype.ValueTranscoder {
	return pgtype.NewEnumType(
		"device_type",
		[]string{
			string(DeviceTypeUndefined),
			string(DeviceTypePhone),
			string(DeviceTypeLaptop),
			string(DeviceTypeIpad),
			string(DeviceTypeDesktop),
			string(DeviceTypeIot),
		},
	)
}

// DeviceType represents the Postgres enum "device_type".
type DeviceType string

const (
	DeviceTypeUndefined DeviceType = "undefined"
	DeviceTypePhone     DeviceType = "phone"
	DeviceTypeLaptop    DeviceType = "laptop"
	DeviceTypeIpad      DeviceType = "ipad"
	DeviceTypeDesktop   DeviceType = "desktop"
	DeviceTypeIot       DeviceType = "iot"
)

func (d DeviceType) String() string { return string(d) }


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

// newDevice creates a new pgtype.ValueTranscoder for the Postgres
// composite type 'device'.
func registernewDevice() pgtype.Codec {
	return tr.newCompositeValue(
		"device",
		compositeField{name: "mac", typeName: "macaddr", defaultCodec: &pgtype.MacaddrCodec{}},
		compositeField{name: "type", typeName: "device_type", defaultCodec: newDeviceTypeEnum()},
	)
}

// newDeviceTypeArray creates a new pgtype.Codec for the Postgres
// '_device_type' array type.
func registernewDeviceTypeArray() pgtype.Codec {
	return tr.newArrayValue("_device_type", "device_type", newDeviceTypeEnum)
}

const findAllDevicesSQL = `SELECT mac, type
FROM device;`

type FindAllDevicesRow struct {
	Mac  pgtype.Macaddr `json:"mac"`
	Type DeviceType     `json:"type"`
}

// FindAllDevices implements Querier.FindAllDevices.
func (q *DBQuerier) FindAllDevices(ctx context.Context) ([]FindAllDevicesRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "FindAllDevices")
	rows, err := q.conn.Query(ctx, findAllDevicesSQL)
	if err != nil {
		return nil, fmt.Errorf("query FindAllDevices: %w", err)
	}

	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (FindAllDevicesRow, error) {
		var item FindAllDevicesRow
		if err := row.Scan(
			&item.Mac, // 'mac', 'Mac', 'pgtype.Macaddr', 'github.com/jackc/pgx/v5/pgtype', 'Macaddr'
			&item.Type, // 'type', 'Type', 'DeviceType', 'github.com/robbert229/pggen/example/enums', 'DeviceType'
		); err != nil {
			return item, fmt.Errorf("failed to scan: %w", err)
		}
		return item, nil
	})
}

const insertDeviceSQL = `INSERT INTO device (mac, type)
VALUES ($1, $2);`

// InsertDevice implements Querier.InsertDevice.
func (q *DBQuerier) InsertDevice(ctx context.Context, mac pgtype.Macaddr, typePg DeviceType) (pgconn.CommandTag, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "InsertDevice")
	cmdTag, err := q.conn.Exec(ctx, insertDeviceSQL, mac, typePg)
	if err != nil {
		return pgconn.CommandTag{}, fmt.Errorf("exec query InsertDevice: %w", err)
	}
	return cmdTag, err
}

const findOneDeviceArraySQL = `SELECT enum_range(NULL::device_type) AS device_types;`

// FindOneDeviceArray implements Querier.FindOneDeviceArray.
func (q *DBQuerier) FindOneDeviceArray(ctx context.Context) ([]DeviceType, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "FindOneDeviceArray")
	rows, err := q.conn.Query(ctx, findOneDeviceArraySQL)
	if err != nil {
		return nil, fmt.Errorf("query FindOneDeviceArray: %w", err)
	}

	return pgx.CollectExactlyOneRow(rows, func(row pgx.CollectableRow) ([]DeviceType, error) {
		var item []DeviceType
		if err := row.Scan(
			&item,
		); err != nil {
			return item, fmt.Errorf("failed to scan: %w", err)
		}
		return item, nil
	})
}

const findManyDeviceArraySQL = `SELECT enum_range('ipad'::device_type, 'iot'::device_type) AS device_types
UNION ALL
SELECT enum_range(NULL::device_type) AS device_types;`

// FindManyDeviceArray implements Querier.FindManyDeviceArray.
func (q *DBQuerier) FindManyDeviceArray(ctx context.Context) ([][]DeviceType, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "FindManyDeviceArray")
	rows, err := q.conn.Query(ctx, findManyDeviceArraySQL)
	if err != nil {
		return nil, fmt.Errorf("query FindManyDeviceArray: %w", err)
	}

	return pgx.CollectRows(rows, func(row pgx.CollectableRow) ([]DeviceType, error) {
		var item []DeviceType
		if err := row.Scan(
			&item,
		); err != nil {
			return item, fmt.Errorf("failed to scan: %w", err)
		}
		return item, nil
	})
}

const findManyDeviceArrayWithNumSQL = `SELECT 1 AS num, enum_range('ipad'::device_type, 'iot'::device_type) AS device_types
UNION ALL
SELECT 2 as num, enum_range(NULL::device_type) AS device_types;`

type FindManyDeviceArrayWithNumRow struct {
	Num         *int32       `json:"num"`
	DeviceTypes []DeviceType `json:"device_types"`
}

// FindManyDeviceArrayWithNum implements Querier.FindManyDeviceArrayWithNum.
func (q *DBQuerier) FindManyDeviceArrayWithNum(ctx context.Context) ([]FindManyDeviceArrayWithNumRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "FindManyDeviceArrayWithNum")
	rows, err := q.conn.Query(ctx, findManyDeviceArrayWithNumSQL)
	if err != nil {
		return nil, fmt.Errorf("query FindManyDeviceArrayWithNum: %w", err)
	}

	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (FindManyDeviceArrayWithNumRow, error) {
		var item FindManyDeviceArrayWithNumRow
		if err := row.Scan(
			&item.Num, // 'num', 'Num', '*int32', '', '*int32'
			&item.DeviceTypes, // 'device_types', 'DeviceTypes', '[]DeviceType', 'github.com/robbert229/pggen/example/enums', '[]DeviceType'
		); err != nil {
			return item, fmt.Errorf("failed to scan: %w", err)
		}
		return item, nil
	})
}

const enumInsideCompositeSQL = `SELECT ROW('08:00:2b:01:02:03'::macaddr, 'phone'::device_type) ::device;`

// EnumInsideComposite implements Querier.EnumInsideComposite.
func (q *DBQuerier) EnumInsideComposite(ctx context.Context) (Device, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "EnumInsideComposite")
	rows, err := q.conn.Query(ctx, enumInsideCompositeSQL)
	if err != nil {
		return Device{}, fmt.Errorf("query EnumInsideComposite: %w", err)
	}

	return pgx.CollectExactlyOneRow(rows, func(row pgx.CollectableRow) (Device, error) {
		var item Device
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