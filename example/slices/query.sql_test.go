package slices

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/robbert229/pggen/internal/difftest"
	"github.com/robbert229/pggen/internal/pgtest"
	"github.com/stretchr/testify/require"
)

func TestNewQuerier_GetBools(t *testing.T) {
	ctx := context.Background()
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

	t.Run("GetBools", func(t *testing.T) {
		want := []bool{true, true, false}
		got, err := q.GetBools(ctx, want)
		require.NoError(t, err)
		difftest.AssertSame(t, want, got)
	})
}

func TestNewQuerier_GetOneTimestamp(t *testing.T) {
	ctx := context.Background()
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
	ts := time.Date(2020, 1, 1, 11, 11, 11, 0, time.UTC)

	t.Run("GetOneTimestamp", func(t *testing.T) {
		got, err := q.GetOneTimestamp(ctx, &ts)
		require.NoError(t, err)
		difftest.AssertSame(t, &ts, got)
	})
}

func TestNewQuerier_GetManyTimestamptzs(t *testing.T) {
	ctx := context.Background()
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
	ts1 := time.Date(2020, 1, 1, 11, 11, 11, 0, time.UTC)
	ts2 := time.Date(2022, 2, 2, 22, 22, 22, 0, time.UTC)

	t.Run("GetManyTimestamptzs", func(t *testing.T) {
		got, err := q.GetManyTimestamptzs(ctx, []time.Time{ts1, ts2})
		require.NoError(t, err)
		difftest.AssertSame(t, []*time.Time{&ts1, &ts2}, got)
	})
}
