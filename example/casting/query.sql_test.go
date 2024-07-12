package casting

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/robbert229/pggen/internal/pgtest"
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

	bookAID := insertBook(t, q, "Book A")
	bookBID := insertBook(t, q, "Book B")

	assignBook(t, q, bookAID, adamsID)
	//assignBook(t, q, bookBID, adamsID)
	_ = bookBID

	results, err := q.FindBooks(context.Background())
	require.NoError(t, err)

	_ = results
}

func insertAuthor(t *testing.T, q *DBQuerier, first, last string) int32 {
	t.Helper()
	authorID, err := q.InsertAuthor(context.Background(), first, last)
	require.NoError(t, err, "insert author")
	return authorID
}

func insertBook(t *testing.T, q *DBQuerier, title string) int32 {
	t.Helper()
	bookID, err := q.InsertBook(context.Background(), title)
	require.NoError(t, err, "insert book")
	return bookID
}

func assignBook(t *testing.T, q *DBQuerier, bookID, authorID int32) {
	t.Helper()
	_, err := q.AssignAuthor(context.Background(), authorID, bookID)
	require.NoError(t, err, "assign author")
}
