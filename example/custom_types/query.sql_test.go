package custom_types

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/robbert229/pggen/internal/pgtest"
	"github.com/robbert229/pggen/internal/texts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQuerier_CustomTypes(t *testing.T) {
	conn, cleanup := pgtest.NewPostgresSchema(t, []string{"schema.sql"})
	defer cleanup()
	q, err := NewQuerier(context.Background(), conn)
	require.NoError(t, err)
	ctx := context.Background()

	t.Run("CustomTypes", func(t *testing.T) {
		val, err := q.CustomTypes(ctx)
		require.NoError(t, err)
		want := CustomTypesRow{
			Column: "some_text",
			Int8:   1,
		}
		assert.Equal(t, want, val)
	})
}

func TestQuerier_CustomMyInt(t *testing.T) {
	conn, cleanup := pgtest.NewPostgresSchema(t, []string{"schema.sql"})
	defer cleanup()
	row := conn.QueryRow(context.Background(), texts.Dedent(`
		SELECT pt.oid
		FROM pg_type pt
			JOIN pg_namespace pn ON pt.typnamespace = pn.oid
		WHERE typname = 'my_int'
			AND pn.nspname = current_schema()
		LIMIT 1;
	`))
	oidVal := uint32(0)
	err := row.Scan(&oidVal)
	require.NoError(t, err)
	t.Logf("my_int oid: %d", oidVal)

	conn.TypeMap().RegisterType(&pgtype.Type{
		Codec: &pgtype.Int2Codec{},
		Name:  "my_int",
		OID:   oidVal,
	})

	q, err := NewQuerier(context.Background(), conn)
	require.NoError(t, err)
	ctx := context.Background()

	t.Run("CustomMyInt", func(t *testing.T) {
		val, err := q.CustomMyInt(ctx)
		require.NoError(t, err)
		assert.Equal(t, 5, val)
	})
}

func TestQuerier_IntArray(t *testing.T) {
	conn, cleanup := pgtest.NewPostgresSchema(t, []string{"schema.sql"})
	defer cleanup()
	q, err := NewQuerier(context.Background(), conn)
	require.NoError(t, err)
	ctx := context.Background()

	t.Run("IntArray", func(t *testing.T) {
		array, err := q.IntArray(ctx)
		require.NoError(t, err)
		assert.Equal(t, [][]int32{{5, 6, 7}}, array)
	})
}
