package pgtype

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/robbert229/pggen/internal/pgtest"
	"github.com/stretchr/testify/assert"
)

func TestNewQuerier_FindAuthorByID(t *testing.T) {
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
	adamsID := insertAuthor(t, q, "john", "adams")
	insertAuthor(t, q, "george", "washington")

	t.Run("FindAuthorByID", func(t *testing.T) {
		authorByID, err := q.FindAuthorByID(context.Background(), pgtype.Int4{Int32: adamsID, Valid: true})
		require.NoError(t, err)
		assert.Equal(t, FindAuthorByIDRow{
			AuthorID:  pgtype.Int4{Int32: adamsID, Valid: true},
			FirstName: pgtype.Text{String: "john", Valid: true},
			LastName:  pgtype.Text{String: "adams", Valid: true},
			Suffix:    pgtype.Text{Valid: false},
		}, authorByID)
	})

	t.Run("FindAuthorByID - none-exists", func(t *testing.T) {
		missingAuthorByID, err := q.FindAuthorByID(context.Background(), pgtype.Int4{Int32: 888, Valid: true})
		require.Error(t, err, "expected error when finding author ID that doesn't match")
		assert.Zero(t, missingAuthorByID, "expected zero value when error")
		if !errors.Is(err, pgx.ErrNoRows) {
			t.Fatalf("expected no rows error to wrap pgx.ErrNoRows; got %s", err)
		}
	})
}

func insertAuthor(t *testing.T, q *DBQuerier, first, last string) int32 {
	t.Helper()
	authorID, err := q.InsertAuthor(context.Background(), pgtype.Text{String: first, Valid: true}, pgtype.Text{String: last, Valid: true})
	require.NoError(t, err, "insert author")
	return authorID.Int32
}
