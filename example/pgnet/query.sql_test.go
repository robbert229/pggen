package pgnet

import (
	"context"
	"errors"
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/robbert229/pggen/internal/pgtest"
	"github.com/stretchr/testify/assert"
)

func TestNewQuerier_FindServerByIP(t *testing.T) {
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

	serverID := insertServer(t, q, &net.IPNet{
		IP: net.IPv4(192, 168, 1, 1),
	})
	insertServer(t, q, &net.IPNet{
		IP: net.IPv4(192, 168, 1, 8),
	})

	t.Run("FindServerByIP", func(t *testing.T) {
		queried, err := q.FindServerByIP(context.Background(), &net.IPNet{
			IP: net.IPv4(192, 168, 1, 1),
		})
		require.NoError(t, err)
		assert.Equal(t, queried.ID, serverID)
		assert.Equal(t, queried.IpAddress.IP, net.IPv4(192, 168, 1, 1))
	})

	t.Run("FindServerByIP - none-exists", func(t *testing.T) {
		missingServerByIP, err := q.FindServerByIP(
			context.Background(),
			&net.IPNet{
				IP: net.IPv4(192, 168, 1, 32),
			},
		)
		require.Error(t, err, "expected error when finding server IP that doesn't match")
		assert.Zero(t, missingServerByIP, "expected zero value when error")
		if !errors.Is(err, pgx.ErrNoRows) {
			t.Fatalf("expected no rows error to wrap pgx.ErrNoRows; got %s", err)
		}
	})
}

func insertServer(t *testing.T, q *DBQuerier, ipAddress *net.IPNet) int32 {
	t.Helper()

	serverID, err := q.InsertServer(context.Background(), ipAddress, nil)
	require.NoError(t, err, "insert server")
	return serverID
}
