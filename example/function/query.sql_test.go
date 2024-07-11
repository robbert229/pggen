package function

import (
	"context"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/robbert229/pggen/internal/difftest"
	"github.com/robbert229/pggen/internal/ptrs"
	"github.com/stretchr/testify/require"

	"github.com/robbert229/pggen/internal/pgtest"
)

func TestNewQuerier_OutParams(t *testing.T) {
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

	t.Run("OutParams", func(t *testing.T) {
		got, err := q.OutParams(context.Background())
		require.NoError(t, err)
		want := []OutParamsRow{
			{
				Items: []*ListItem{{Name: ptrs.String("some_name"), Color: ptrs.String("some_color")}},
				Stats: &ListStats{
					Val1: ptrs.String("abc"),
					Val2: []*int32{ptrs.Int32(1), ptrs.Int32(2)},
				},
			},
		}
		difftest.AssertSame(t, want, got)
	})
}
