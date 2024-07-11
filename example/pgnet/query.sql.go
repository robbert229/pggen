// Code generated by pggen. DO NOT EDIT.

package pgnet

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
	// FindServers finds all servers.
	FindServers(ctx context.Context) ([]FindServersRow, error)

	// FindServerByIP finds a server by its ip address.
	FindServerByIP(ctx context.Context, ipAddress *net.IPNet) (FindServerByIPRow, error)

	// InsertServer inserts a server and returns the ID.
	InsertServer(ctx context.Context, ipAddress *net.IPNet, extraIpAddress *net.IPNet) (int32, error)
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



const findServersSQL = `SELECT * FROM servers;`

type FindServersRow struct {
	ID             *int32     `json:"id"`
	IpAddress      *net.IPNet `json:"ip_address"`
	ExtraIpAddress *net.IPNet `json:"extra_ip_address"`
}

// FindServers implements Querier.FindServers.
func (q *DBQuerier) FindServers(ctx context.Context) ([]FindServersRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "FindServers")
	rows, err := q.conn.Query(ctx, findServersSQL)
	if err != nil {
		return nil, fmt.Errorf("query FindServers: %w", err)
	}

	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (FindServersRow, error) {
  var item FindServersRow
		if err := row.Scan(&item.ID, // 'id', 'ID', '*int32', '', '*int32'
			&item.IpAddress, // 'ip_address', 'IpAddress', '*net.IPNet', '', '*IPNet'
			&item.ExtraIpAddress, // 'extra_ip_address', 'ExtraIpAddress', '*net.IPNet', '', '*IPNet'
			); err != nil {
			return item, fmt.Errorf("failed to scan: %w", err)
		}
		return item, nil
	})
}

const findServerByIPSQL = `SELECT * FROM servers WHERE ip_address = $1;`

type FindServerByIPRow struct {
	ID             int32      `json:"id"`
	IpAddress      *net.IPNet `json:"ip_address"`
	ExtraIpAddress *net.IPNet `json:"extra_ip_address"`
}

// FindServerByIP implements Querier.FindServerByIP.
func (q *DBQuerier) FindServerByIP(ctx context.Context, ipAddress *net.IPNet) (FindServerByIPRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "FindServerByIP")
	rows, err := q.conn.Query(ctx, findServerByIPSQL, ipAddress)
	if err != nil {
		return FindServerByIPRow{}, fmt.Errorf("query FindServerByIP: %w", err)
	}

	return pgx.CollectExactlyOneRow(rows, func(row pgx.CollectableRow) (FindServerByIPRow, error) {
  var item FindServerByIPRow
		if err := row.Scan(&item.ID, // 'id', 'ID', 'int32', '', 'int32'
			&item.IpAddress, // 'ip_address', 'IpAddress', '*net.IPNet', '', '*IPNet'
			&item.ExtraIpAddress, // 'extra_ip_address', 'ExtraIpAddress', '*net.IPNet', '', '*IPNet'
			); err != nil {
			return item, fmt.Errorf("failed to scan: %w", err)
		}
		return item, nil
	})
}

const insertServerSQL = `INSERT INTO  servers (ip_address, extra_ip_address)
VALUES ($1, $2)
RETURNING id;`

// InsertServer implements Querier.InsertServer.
func (q *DBQuerier) InsertServer(ctx context.Context, ipAddress *net.IPNet, extraIpAddress *net.IPNet) (int32, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "InsertServer")
	rows, err := q.conn.Query(ctx, insertServerSQL, ipAddress, extraIpAddress)
	if err != nil {
		return 0, fmt.Errorf("query InsertServer: %w", err)
	}

	return pgx.CollectExactlyOneRow(rows, func(row pgx.CollectableRow) (int32, error) {
  var item int32
		if err := row.Scan(&item,
			); err != nil {
			return item, fmt.Errorf("failed to scan: %w", err)
		}
		return item, nil
	})
}
