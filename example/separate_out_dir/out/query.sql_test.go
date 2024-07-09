package out

import (
	"context"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"

	"github.com/robbert229/pggen/internal/pgtest"
	"github.com/stretchr/testify/assert"
)

func TestNewQuerier_FindAuthorByID(t *testing.T) {
	pool, cleanup := pgtest.NewPostgresSchema(t, []string{"../schema.sql"}, func(config *pgxpool.Config) {
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

	t.Run("AlphaNested", func(t *testing.T) {
		got, err := q.AlphaNested(context.Background())
		require.NoError(t, err)
		assert.Equal(t, "alpha_nested", got)
	})

	t.Run("Alpha", func(t *testing.T) {
		got, err := q.Alpha(context.Background())
		require.NoError(t, err)
		assert.Equal(t, "alpha", got)
	})

	t.Run("Bravo", func(t *testing.T) {
		got, err := q.Bravo(context.Background())
		require.NoError(t, err)
		assert.Equal(t, "bravo", got)
	})
}
