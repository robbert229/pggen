package complex_params

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

func TestNewQuerier_ParamArrayInt(t *testing.T) {
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

	want := []int{1, 2, 3, 4}

	t.Run("ParamArrayInt", func(t *testing.T) {
		row, err := q.ParamArrayInt(ctx, want)
		require.NoError(t, err)
		assert.Equal(t, want, row)
	})
}

func TestNewQuerier_ParamNested1(t *testing.T) {
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

	want := Dimensions{Width: 77, Height: 77}

	t.Run("ParamNested1", func(t *testing.T) {
		row, err := q.ParamNested1(ctx, want)
		require.NoError(t, err)
		assert.Equal(t, want, row)
	})
}

func TestNewQuerier_ParamNested2(t *testing.T) {
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

	want := ProductImageType{
		Source:     "src",
		Dimensions: &Dimensions{Width: 77, Height: 77},
	}

	t.Run("ParamNested2", func(t *testing.T) {
		row, err := q.ParamNested2(ctx, want)
		require.NoError(t, err)
		assert.Equal(t, want, row)
	})
}

func TestNewQuerier_ParamNested2Array(t *testing.T) {
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
		{Source: "src1", Dimensions: &Dimensions{Width: 11, Height: 11}},
		{Source: "src2", Dimensions: &Dimensions{Width: 22, Height: 22}},
	}

	t.Run("ParamNested2Array", func(t *testing.T) {
		row, err := q.ParamNested2Array(ctx, want)
		require.NoError(t, err)
		assert.Equal(t, want, row)
	})
}

func TestNewQuerier_ParamNested3(t *testing.T) {
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

	want := ProductImageSetType{
		Name:      "set1",
		OrigImage: &ProductImageType{Source: "src1", Dimensions: &Dimensions{Width: 11, Height: 11}},
		Images: []*ProductImageType{
			{Source: "src1", Dimensions: &Dimensions{Width: 11, Height: 11}},
			{Source: "src2", Dimensions: &Dimensions{Width: 22, Height: 22}},
		},
	}

	t.Run("ParamNested3", func(t *testing.T) {
		row, err := q.ParamNested3(ctx, want)
		require.NoError(t, err)
		assert.Equal(t, want, row)
	})
}

func TestNewQuerier_ParamNested3_QueryAllDataTypes(t *testing.T) {
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
	ctx := context.Background()
	// dataTypes, err := QueryAllDataTypes(ctx, conn)
	// require.NoError(t, err)
	conn, err := pool.Acquire(context.Background())
	require.NoError(t, err)
	defer conn.Release()

	q, err := NewQuerier(context.Background(), conn)

	require.NoError(t, err)

	want := ProductImageSetType{
		Name:      "set1",
		OrigImage: &ProductImageType{Source: "src1", Dimensions: &Dimensions{Width: 11, Height: 11}},
		Images: []*ProductImageType{
			{Source: "src1", Dimensions: &Dimensions{Width: 11, Height: 11}},
			{Source: "src2", Dimensions: &Dimensions{Width: 22, Height: 22}},
		},
	}

	t.Run("ParamNested3", func(t *testing.T) {
		row, err := q.ParamNested3(ctx, want)
		require.NoError(t, err)
		assert.Equal(t, want, row)
	})
}
