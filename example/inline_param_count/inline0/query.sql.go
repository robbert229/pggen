// Code generated by pggen. DO NOT EDIT.

package inline0

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

var _ genericConn = (*pgx.Conn)(nil)
var _ RegisterConn = (*pgx.Conn)(nil)

// Querier is a typesafe Go interface backed by SQL queries.
type Querier interface {
	// CountAuthors returns the number of authors (zero params).
	CountAuthors(ctx context.Context) (*int, error)

	// FindAuthorById finds one (or zero) authors by ID (one param).
	FindAuthorByID(ctx context.Context, params FindAuthorByIDParams) (FindAuthorByIDRow, error)

	// InsertAuthor inserts an author by name and returns the ID (two params).
	InsertAuthor(ctx context.Context, params InsertAuthorParams) (int32, error)

	// DeleteAuthorsByFullName deletes authors by the full name (three params).
	DeleteAuthorsByFullName(ctx context.Context, params DeleteAuthorsByFullNameParams) (pgconn.CommandTag, error)
}

var _ Querier = &DBQuerier{}

type DBQuerier struct {
	conn  genericConn
}

// genericConn is a connection like *pgx.Conn, pgx.Tx, or *pgxpool.Pool.
type genericConn interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
}

// NewQuerier creates a DBQuerier that implements Querier.
func NewQuerier(ctx context.Context, conn genericConn) (*DBQuerier, error) {
	return &DBQuerier{
		conn: conn, 
	}, nil
}

type typeHook func(ctx context.Context, conn RegisterConn) error

var typeHooks []typeHook

func addHook(hook typeHook) {
	typeHooks = append(typeHooks, hook)
}

type RegisterConn interface {
	LoadType(ctx context.Context, typeName string) (*pgtype.Type, error)
	TypeMap() *pgtype.Map
}

func Register(ctx context.Context, conn RegisterConn) error {
  

	for _, hook := range typeHooks {
		if err := hook(ctx, conn); err != nil {
			return err
		}
	}
	
	return nil
}



const countAuthorsSQL = `SELECT count(*) FROM author;`

// CountAuthors implements Querier.CountAuthors.
func (q *DBQuerier) CountAuthors(ctx context.Context) (*int, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "CountAuthors")
	rows, err := q.conn.Query(ctx, countAuthorsSQL)
	if err != nil {
		return nil, fmt.Errorf("query CountAuthors: %w", err)
	}

	return pgx.CollectExactlyOneRow(rows, func(row pgx.CollectableRow) (*int, error) {
  var item *int
		if err := row.Scan(&item,
			); err != nil {
			return item, fmt.Errorf("failed to scan: %w", err)
		}
		return item, nil
	})
}

const findAuthorByIDSQL = `SELECT * FROM author WHERE author_id = $1;`

type FindAuthorByIDParams struct {
	AuthorID int32 `json:"AuthorID"`
}

type FindAuthorByIDRow struct {
	AuthorID  int32   `json:"author_id"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Suffix    *string `json:"suffix"`
}

// FindAuthorByID implements Querier.FindAuthorByID.
func (q *DBQuerier) FindAuthorByID(ctx context.Context, params FindAuthorByIDParams) (FindAuthorByIDRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "FindAuthorByID")
	rows, err := q.conn.Query(ctx, findAuthorByIDSQL, params.AuthorID)
	if err != nil {
		return FindAuthorByIDRow{}, fmt.Errorf("query FindAuthorByID: %w", err)
	}

	return pgx.CollectExactlyOneRow(rows, func(row pgx.CollectableRow) (FindAuthorByIDRow, error) {
  var item FindAuthorByIDRow
		if err := row.Scan(&item.AuthorID, // 'author_id', 'AuthorID', 'int32', '', 'int32'
			&item.FirstName, // 'first_name', 'FirstName', 'string', '', 'string'
			&item.LastName, // 'last_name', 'LastName', 'string', '', 'string'
			&item.Suffix, // 'suffix', 'Suffix', '*string', '', '*string'
			); err != nil {
			return item, fmt.Errorf("failed to scan: %w", err)
		}
		return item, nil
	})
}

const insertAuthorSQL = `INSERT INTO author (first_name, last_name)
VALUES ($1, $2)
RETURNING author_id;`

type InsertAuthorParams struct {
	FirstName string `json:"FirstName"`
	LastName  string `json:"LastName"`
}

// InsertAuthor implements Querier.InsertAuthor.
func (q *DBQuerier) InsertAuthor(ctx context.Context, params InsertAuthorParams) (int32, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "InsertAuthor")
	rows, err := q.conn.Query(ctx, insertAuthorSQL, params.FirstName, params.LastName)
	if err != nil {
		return 0, fmt.Errorf("query InsertAuthor: %w", err)
	}

	return pgx.CollectExactlyOneRow(rows, func(row pgx.CollectableRow) (int32, error) {
  var item int32
		if err := row.Scan(&item,
			); err != nil {
			return item, fmt.Errorf("failed to scan: %w", err)
		}
		return item, nil
	})
}

const deleteAuthorsByFullNameSQL = `DELETE
FROM author
WHERE first_name = $1
  AND last_name = $2
  AND CASE WHEN $3 = '' THEN suffix IS NULL ELSE suffix = $3 END;`

type DeleteAuthorsByFullNameParams struct {
	FirstName string `json:"FirstName"`
	LastName  string `json:"LastName"`
	Suffix    string `json:"Suffix"`
}

// DeleteAuthorsByFullName implements Querier.DeleteAuthorsByFullName.
func (q *DBQuerier) DeleteAuthorsByFullName(ctx context.Context, params DeleteAuthorsByFullNameParams) (pgconn.CommandTag, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "DeleteAuthorsByFullName")
	cmdTag, err := q.conn.Exec(ctx, deleteAuthorsByFullNameSQL, params.FirstName, params.LastName, params.Suffix)
	if err != nil {
		return pgconn.CommandTag{}, fmt.Errorf("exec query DeleteAuthorsByFullName: %w", err)
	}
	return cmdTag, err
}
