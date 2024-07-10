package ltree

import (
	"context"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
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

	if _, err := q.InsertSampleData(ctx); err != nil {
		t.Fatal(err)
	}

	{
		rows, err := q.FindTopScienceChildren(ctx)
		require.NoError(t, err)
		want := []pgtype.Text{
			{String: "Top.Science", Valid: true},
			{String: "Top.Science.Astronomy", Valid: true},
			{String: "Top.Science.Astronomy.Astrophysics", Valid: true},
			{String: "Top.Science.Astronomy.Cosmology", Valid: true},
		}
		assert.Equal(t, want, rows)
	}

	{
		rows, err := q.FindTopScienceChildrenAgg(ctx)
		require.NoError(t, err)
		want := []pgtype.Text{
			{String: "Top.Science", Valid: true},
			{String: "Top.Science.Astronomy", Valid: true},
			{String: "Top.Science.Astronomy.Astrophysics", Valid: true},
			{String: "Top.Science.Astronomy.Cosmology", Valid: true},
		}
		assert.Equal(t, want, rows)
	}

	{
		in1 := pgtype.Text{String: "foo", Valid: true}
		in2 := []string{"qux", "qux"}
		in2Txt := newTextArray(in2)
		rows, err := q.FindLtreeInput(ctx, in1, in2)
		require.NoError(t, err)
		assert.Equal(t, FindLtreeInputRow{
			Ltree:   in1,
			TextArr: in2Txt,
		}, rows)
	}
}

// newTextArray creates a one dimensional text array from the string slice with
// no null elements.
func newTextArray(ss []string) []pgtype.Text {
	elems := make([]pgtype.Text, len(ss))
	for i, s := range ss {
		elems[i] = pgtype.Text{String: s, Valid: true}
	}

	return elems
}
