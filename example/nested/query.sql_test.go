package nested

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

func TestNewQuerier_ArrayNested2(t *testing.T) {
	t.SkipNow()

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

	want := []ProductImageType{
		{Source: "img2", Dimensions: Dimensions{22, 22}},
		{Source: "img3", Dimensions: Dimensions{33, 33}},
	}
	t.Run("ArrayNested2", func(t *testing.T) {
		rows, err := q.ArrayNested2(ctx)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, want, rows)
	})
}

func TestNewQuerier_Nested3(t *testing.T) {
	t.SkipNow()

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

	want := []ProductImageSetType{
		{
			Name: "name",
			OrigImage: ProductImageType{
				Source:     "img1",
				Dimensions: Dimensions{Width: 11, Height: 11},
			},
			Images: []ProductImageType{
				{Source: "img2", Dimensions: Dimensions{22, 22}},
				{Source: "img3", Dimensions: Dimensions{33, 33}},
			},
		},
	}
	t.Run("Nested3", func(t *testing.T) {
		t.Skipf("https://github.com/jackc/pgx/issues/874")
		rows, err := q.Nested3(ctx)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, want, rows)
	})
}
