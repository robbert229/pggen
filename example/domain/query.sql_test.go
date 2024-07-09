package domain

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

func TestQuerier_DomainOne(t *testing.T) {
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

	t.Run("DomainOne", func(t *testing.T) {
		postCode, err := q.DomainOne(ctx)
		require.NoError(t, err)
		assert.Equal(t, "90210", postCode)
	})
}
