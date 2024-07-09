package void

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

func TestQuerier(t *testing.T) {
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

	_, err = q.VoidOnly(ctx)
	require.NoError(t, err)

	_, err = q.VoidOnlyTwoParams(ctx, 33)
	require.NoError(t, err)

	{
		row, err := q.VoidTwo(ctx)
		require.NoError(t, err)

		assert.Equal(t, "foo", row)
	}

	{
		row, err := q.VoidThree(ctx)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, VoidThreeRow{Foo: "foo", Bar: "bar"}, row)
	}

	{
		foos, err := q.VoidThree2(ctx)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, []string{"foo"}, foos)
	}
}
