// Code generated by pggen. DO NOT EDIT.

package enums

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"net"
)

var _ genericConn = (*pgx.Conn)(nil)
var _ RegisterConn = (*pgx.Conn)(nil)

// Querier is a typesafe Go interface backed by SQL queries.
type Querier interface {
	FindAllDevices(ctx context.Context) ([]FindAllDevicesRow, error)

	InsertDevice(ctx context.Context, mac net.HardwareAddr, typePg DeviceType) (pgconn.CommandTag, error)

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



// Device represents the Postgres composite type "device".
type Device struct {
	Mac  net.HardwareAddr `json:"mac"`
	Type DeviceType       `json:"type"`
}


// register_newDeviceTypeEnum registers the given postgres type with pgx.
func register_newDeviceTypeEnum(
	ctx context.Context,
	conn RegisterConn,
) error {
	t, err := conn.LoadType(
		ctx,
		"\"device_type\"",
	)
	if err != nil {
		return fmt.Errorf("newDeviceTypeEnum failed to load type: %w", err)
	}
	
	conn.TypeMap().RegisterType(t)
	
	t, err = conn.LoadType(
		ctx,
		"_device_type",
	)
	if err != nil {
		return fmt.Errorf("newDeviceTypeEnum failed to load type: %w", err)
	}
	
	conn.TypeMap().RegisterType(t)
	
	return nil
}

func codec_newDeviceTypeEnum(conn RegisterConn) (pgtype.Codec, error) {
	return &pgtype.EnumCodec{}, nil
}

func init(){
	addHook(register_newDeviceTypeEnum) 
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




	// codec_newDevice is a codec for the composite type of the same name
	func codec_newDevice(conn RegisterConn) (pgtype.Codec, error) {
		
		    field0, ok := conn.TypeMap().TypeForName("macaddr")
			if !ok {
				return nil, fmt.Errorf("type not found: macaddr")
			}
		
		    field1, ok := conn.TypeMap().TypeForName("device_type")
			if !ok {
				return nil, fmt.Errorf("type not found: device_type")
			}
		
		
		return &pgtype.CompositeCodec{
			Fields: []pgtype.CompositeCodecField{
				
					{
						Name: "mac",
						Type: field0,
					},
				
					{
						Name: "type",
						Type: field1,
					},
				
			},
		}, nil
	}

	func register_newDevice(
		ctx context.Context,
		conn RegisterConn,
	) error {
		t, err := conn.LoadType(
			ctx,
			"\"device\"",
		)
		if err != nil {
			return fmt.Errorf("newDevice failed to load type: %w", err)
		}
		
		conn.TypeMap().RegisterType(t)

		return nil
	}

	func init(){
		addHook(register_newDevice) 
	}
	


	// codec_newDeviceTypeArray is a codec for the composite type of the same name
	func codec_newDeviceTypeArray(conn RegisterConn) (pgtype.Codec, error) {
		elementType, ok := conn.TypeMap().TypeForName("device_type")
		if !ok {
			return nil, fmt.Errorf("type not found: device_type")
		}

		return &pgtype.ArrayCodec{
			ElementType: elementType,
		}, nil
	}

	func register_newDeviceTypeArray(
		ctx context.Context,
		conn RegisterConn,
	) error {
		t, err := conn.LoadType(
			ctx,
			"\"_device_type\"",
		)
		if err != nil {
			return fmt.Errorf("newDeviceTypeArray failed to load type: %w", err)
		}

		conn.TypeMap().RegisterType(t)

		return nil
	}

	func init(){
		addHook(register_newDeviceTypeArray) 
	}
	

const findAllDevicesSQL = `SELECT mac, type
FROM device;`

type FindAllDevicesRow struct {
	Mac  net.HardwareAddr `json:"mac"`
	Type DeviceType       `json:"type"`
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
		if err := row.Scan(&item.Mac, // 'mac', 'Mac', 'net.HardwareAddr', 'net', 'HardwareAddr'
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
func (q *DBQuerier) InsertDevice(ctx context.Context, mac net.HardwareAddr, typePg DeviceType) (pgconn.CommandTag, error) {
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
		if err := row.Scan(&item,
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
		if err := row.Scan(&item,
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
		if err := row.Scan(&item.Num, // 'num', 'Num', '*int32', '', '*int32'
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
		if err := row.Scan(&item,
			); err != nil {
			return item, fmt.Errorf("failed to scan: %w", err)
		}
		return item, nil
	})
}
