package device

import (
	"context"
	"fmt"
	"net"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/robbert229/pggen/internal/pgtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQuerier_FindDevicesByUser(t *testing.T) {
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
	ctx := context.Background()
	userID := 18

	_, err = q.InsertUser(ctx, userID, "foo")
	require.NoError(t, err)
	mac1, _ := net.ParseMAC("11:22:33:44:55:66")
	_, err = q.InsertDevice(ctx, mac1, userID)
	require.NoError(t, err)

	t.Run("FindDevicesByUser", func(t *testing.T) {
		val, err := q.FindDevicesByUser(ctx, userID)
		require.NoError(t, err)
		want := []FindDevicesByUserRow{
			{
				ID:   userID,
				Name: "foo",
				MacAddrs: []net.HardwareAddr{
					mac1,
				},
			},
		}
		assert.Equal(t, want, val)
	})
}

var _ genericConn = (*pgx.Conn)(nil)

func TestQuerier_CompositeUser(t *testing.T) {
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
	ctx := context.Background()

	userID := 18
	name := "foo"
	_, err = q.InsertUser(ctx, userID, name)
	require.NoError(t, err)

	mac1, _ := net.ParseMAC("11:22:33:44:55:66")
	mac2, _ := net.ParseMAC("aa:bb:cc:dd:ee:ff")
	_, err = q.InsertDevice(ctx, mac1, userID)
	require.NoError(t, err)
	_, err = q.InsertDevice(ctx, mac2, userID)
	require.NoError(t, err)

	t.Run("CompositeUser", func(t *testing.T) {
		users, err := q.CompositeUser(ctx)
		require.NoError(t, err)
		want := []CompositeUserRow{
			{
				Mac:  mac1,
				Type: DeviceTypeUndefined,
				User: &User{ID: &userID, Name: &name},
			},
			{
				Mac:  mac2,
				Type: DeviceTypeUndefined,
				User: &User{ID: &userID, Name: &name},
			},
		}
		assert.Equal(t, want, users)
	})
}

func TestQuerier_CompositeUserOne(t *testing.T) {
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
	ctx := context.Background()
	id := 15
	name := "qux"
	wantUser := User{ID: &id, Name: &name}

	t.Run("CompositeUserOne", func(t *testing.T) {
		got, err := q.CompositeUserOne(ctx)
		require.NoError(t, err)
		assert.Equal(t, wantUser, got)
	})
}
