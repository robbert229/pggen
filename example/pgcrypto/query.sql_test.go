package pgcrypto

import (
	"context"
	"strings"
	"testing"

	"github.com/robbert229/pggen/internal/pgtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQuerier(t *testing.T) {
	conn, cleanup := pgtest.NewPostgresSchema(t, []string{"schema.sql"})
	defer cleanup()

	q, err := NewQuerier(context.Background(), conn)
	require.NoError(t, err)
	ctx := context.Background()

	_, err = q.CreateUser(ctx, "foo", "hunter2")
	if err != nil {
		t.Fatal(err)
	}

	row, err := q.FindUser(ctx, "foo")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "foo", row.Email, "email should match")
	if !strings.HasPrefix(row.Pass, "$2a$") {
		t.Fatalf("expected hashed password to have prefix $2a$; got %s", row.Pass)
	}
}
