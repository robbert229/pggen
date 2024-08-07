package syntax

import (
	"context"
	"testing"

	"github.com/robbert229/pggen/internal/pgtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQuerier(t *testing.T) {
	pool, cleanup := pgtest.NewPostgresSchema(t, nil)
	defer cleanup()

	conn, err := pool.Acquire(context.Background())
	require.NoError(t, err)
	defer conn.Release()

	q, err := NewQuerier(context.Background(), conn)

	require.NoError(t, err)
	ctx := context.Background()

	val, err := q.Backtick(ctx)
	assert.NoError(t, err, "Backtick")
	assert.Equal(t, "`", val, "Backtick")

	val, err = q.BacktickDoubleQuote(ctx)
	assert.NoError(t, err, "BacktickDoubleQuote")
	assert.Equal(t, "`\"", val, "BacktickDoubleQuote")

	val, err = q.BacktickQuoteBacktick(ctx)
	assert.NoError(t, err, "BacktickQuoteBacktick")
	assert.Equal(t, "`\"`", val, "BacktickQuoteBacktick")

	val, err = q.BacktickNewline(ctx)
	assert.NoError(t, err, "BacktickNewline")
	assert.Equal(t, "`\n", val, "BacktickNewline")

	val, err = q.BacktickBackslashN(ctx)
	assert.NoError(t, err, "BacktickBackslashN")
	assert.Equal(t, "`\\n", val, "BacktickBackslashN")
}
