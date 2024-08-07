package enums

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/robbert229/pggen/internal/pgtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewQuerier_FindAllDevices(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	pool, cleanup := pgtest.NewPostgresSchema(t, []string{"schema.sql"}, func(config *pgxpool.Config) {
		config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
			err := Register(ctx, conn)
			if err != nil {
				return fmt.Errorf("failed to register types: %w", err)
			}

			return nil
		}
	})

	defer cleanup()

	conn, err := pool.Acquire(context.Background())
	require.NoError(t, err)
	defer conn.Release()

	q, err := NewQuerier(context.Background(), conn)

	require.NoError(t, err)

	mac, _ := net.ParseMAC("00:00:5e:00:53:01")

	insertDevice(t, q, mac, DeviceTypeIot)

	t.Run("FindAllDevices", func(t *testing.T) {
		devices, err := q.FindAllDevices(ctx)
		require.NoError(t, err)
		assert.Equal(t,
			[]FindAllDevicesRow{
				{Mac: mac, Type: DeviceTypeIot},
			},
			devices,
		)
	})
}

var allDeviceTypes = []DeviceType{
	DeviceTypeUndefined,
	DeviceTypePhone,
	DeviceTypeLaptop,
	DeviceTypeIpad,
	DeviceTypeDesktop,
	DeviceTypeIot,
}

func TestNewQuerier_FindOneDeviceArray(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	pool, cleanup := pgtest.NewPostgresSchema(t, []string{"schema.sql"}, func(config *pgxpool.Config) {
		config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
			err := Register(ctx, conn)
			if err != nil {
				return fmt.Errorf("failed to register types: %w", err)
			}

			return nil
		}
	})

	defer cleanup()

	conn, err := pool.Acquire(context.Background())
	require.NoError(t, err)
	defer conn.Release()

	q, err := NewQuerier(context.Background(), conn)

	require.NoError(t, err)

	t.Run("FindOneDeviceArray", func(t *testing.T) {
		devices, err := q.FindOneDeviceArray(ctx)
		require.NoError(t, err)
		assert.Equal(t, allDeviceTypes, devices)
	})
}

func TestNewQuerier_FindManyDeviceArray(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	pool, cleanup := pgtest.NewPostgresSchema(t, []string{"schema.sql"}, func(config *pgxpool.Config) {
		config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
			err := Register(ctx, conn)
			if err != nil {
				return fmt.Errorf("failed to register types: %w", err)
			}

			return nil
		}
	})

	defer cleanup()

	conn, err := pool.Acquire(context.Background())
	require.NoError(t, err)
	defer conn.Release()

	q, err := NewQuerier(context.Background(), conn)

	require.NoError(t, err)

	t.Run("FindManyDeviceArray", func(t *testing.T) {
		devices, err := q.FindManyDeviceArray(ctx)
		require.NoError(t, err)
		assert.Equal(t, [][]DeviceType{allDeviceTypes[3:], allDeviceTypes}, devices)
	})
}

func TestNewQuerier_FindManyDeviceArrayWithNum(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	pool, cleanup := pgtest.NewPostgresSchema(t, []string{"schema.sql"}, func(config *pgxpool.Config) {
		config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
			err := Register(ctx, conn)
			if err != nil {
				return fmt.Errorf("failed to register types: %w", err)
			}

			return nil
		}
	})

	defer cleanup()

	conn, err := pool.Acquire(context.Background())
	require.NoError(t, err)
	defer conn.Release()

	q, err := NewQuerier(context.Background(), conn)

	require.NoError(t, err)

	one, two := int32(1), int32(2)

	t.Run("FindManyDeviceArrayWithNum", func(t *testing.T) {
		devices, err := q.FindManyDeviceArrayWithNum(ctx)
		require.NoError(t, err)
		assert.Equal(t, []FindManyDeviceArrayWithNumRow{
			{Num: &one, DeviceTypes: allDeviceTypes[3:]},
			{Num: &two, DeviceTypes: allDeviceTypes},
		}, devices)
	})
}

func TestNewQuerier_EnumInsideComposite(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	pool, cleanup := pgtest.NewPostgresSchema(t, []string{"schema.sql"}, func(config *pgxpool.Config) {
		config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
			err := Register(ctx, conn)
			if err != nil {
				return fmt.Errorf("failed to register types: %w", err)
			}

			return nil
		}
	})

	defer cleanup()

	conn, err := pool.Acquire(context.Background())
	require.NoError(t, err)
	defer conn.Release()

	q, err := NewQuerier(context.Background(), conn)

	require.NoError(t, err)

	mac, _ := net.ParseMAC("08:00:2b:01:02:03")

	t.Run("EnumInsideComposite", func(t *testing.T) {
		device, err := q.EnumInsideComposite(ctx)
		require.NoError(t, err)
		assert.Equal(t,
			Device{Mac: mac, Type: DeviceTypePhone},
			device,
		)
	})
}

func insertDevice(t *testing.T, q *DBQuerier, mac net.HardwareAddr, device DeviceType) {
	t.Helper()
	_, err := q.InsertDevice(context.Background(),
		mac,
		device,
	)
	require.NoError(t, err)
}
