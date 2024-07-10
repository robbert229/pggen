package go_pointer_types

import (
	"context"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/robbert229/pggen/internal/pgtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQuerier_GenSeries1(t *testing.T) {
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

	t.Run("GenSeries1", func(t *testing.T) {
		got, err := q.GenSeries1(ctx)
		require.NoError(t, err)
		zero := 0
		assert.Equal(t, &zero, got)
	})
}

func TestQuerier_GenSeries(t *testing.T) {
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

	t.Run("GenSeries", func(t *testing.T) {
		got, err := q.GenSeries(ctx)
		if err != nil {
			t.Fatal(err)
		}
		zero, one, two := 0, 1, 2
		assert.Equal(t, []*int{&zero, &one, &two}, got)
	})
}

func TestQuerier_GenSeriesArr1(t *testing.T) {
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

	t.Run("GenSeriesArr1", func(t *testing.T) {
		got, err := q.GenSeriesArr1(ctx)
		require.NoError(t, err)
		assert.Equal(t, []int{0, 1, 2}, got)
	})
}

func TestQuerier_GenSeriesArr(t *testing.T) {
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

	t.Run("GenSeriesArr", func(t *testing.T) {
		got, err := q.GenSeriesArr(ctx)
		require.NoError(t, err)
		assert.Equal(t, [][]int{{0, 1, 2}}, got)
	})
}

func TestQuerier_GenSeriesStr(t *testing.T) {
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

	t.Run("GenSeriesStr1", func(t *testing.T) {
		got, err := q.GenSeriesStr1(ctx)
		require.NoError(t, err)
		zero := "0"
		assert.Equal(t, &zero, got)
	})

	t.Run("GenSeriesStr", func(t *testing.T) {
		got, err := q.GenSeriesStr(ctx)
		require.NoError(t, err)
		zero, one, two := "0", "1", "2"
		assert.Equal(t, []*string{&zero, &one, &two}, got)
	})
}
