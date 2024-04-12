package out

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/robbert229/pggen/internal/pgtest"
	"github.com/stretchr/testify/assert"
)

func TestNewQuerier_FindAuthorByID(t *testing.T) {
	conn, cleanup := pgtest.NewPostgresSchema(t, []string{"../schema.sql"})
	defer cleanup()

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
