package pgcrypto

import (
	"context"
	"fmt"
	"strings"
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
