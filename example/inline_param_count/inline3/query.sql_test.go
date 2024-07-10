package inline3

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/robbert229/pggen/internal/pgtest"
	"github.com/stretchr/testify/assert"
)

func TestNewQuerier_FindAuthorByID(t *testing.T) {
	pool, cleanup := pgtest.NewPostgresSchema(t, []string{"../schema.sql"}, func(config *pgxpool.Config) {
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

	t.Run("CountAuthors two", func(t *testing.T) {
		got, err := q.CountAuthors(context.Background())
		require.NoError(t, err)
		assert.Equal(t, 2, *got)
	})

	t.Run("FindAuthorByID", func(t *testing.T) {
		authorByID, err := q.FindAuthorByID(context.Background(), adamsID)
		require.NoError(t, err)
		assert.Equal(t, FindAuthorByIDRow{
			AuthorID:  adamsID,
			FirstName: "john",
			LastName:  "adams",
			Suffix:    nil,
		}, authorByID)
	})

	t.Run("FindAuthorByID - none-exists", func(t *testing.T) {
		missingAuthorByID, err := q.FindAuthorByID(context.Background(), 888)
		require.Error(t, err, "expected error when finding author ID that doesn't match")
		assert.Zero(t, missingAuthorByID, "expected zero value when error")
		if !errors.Is(err, pgx.ErrNoRows) {
			t.Fatalf("expected no rows error to wrap pgx.ErrNoRows; got %s", err)
		}
	})
}

func TestNewQuerier_DeleteAuthorsByFullName(t *testing.T) {
	pool, cleanup := pgtest.NewPostgresSchema(t, []string{"../schema.sql"}, func(config *pgxpool.Config) {
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
	insertAuthor(t, q, "george", "washington")

	t.Run("DeleteAuthorsByFullName", func(t *testing.T) {
		tag, err := q.DeleteAuthorsByFullName(context.Background(), "george", "washington", "")
		require.NoError(t, err)
		assert.Truef(t, tag.Delete(), "expected delete tag; got %s", tag.String())
		assert.Equal(t, int64(1), tag.RowsAffected())
	})
}

func insertAuthor(t *testing.T, q *DBQuerier, first, last string) int32 {
	t.Helper()
	authorID, err := q.InsertAuthor(context.Background(), first, last)
	require.NoError(t, err, "insert author")
	return authorID
}
