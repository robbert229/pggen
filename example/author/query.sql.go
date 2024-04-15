// Code generated by pggen. DO NOT EDIT.

package author

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

var _ genericConn = (*pgx.Conn)(nil)

// Querier is a typesafe Go interface backed by SQL queries.
type Querier interface {
	// FindAuthorById finds one (or zero) authors by ID.
	FindAuthorByID(ctx context.Context, authorID int32) (FindAuthorByIDRow, error)

	// FindAuthors finds authors by first name.
	FindAuthors(ctx context.Context, firstName string) ([]FindAuthorsRow, error)

	// FindAuthorNames finds one (or zero) authors by ID.
	FindAuthorNames(ctx context.Context, authorID int32) ([]FindAuthorNamesRow, error)

	// FindFirstNames finds one (or zero) authors by ID.
	FindFirstNames(ctx context.Context, authorID int32) ([]*string, error)

	// DeleteAuthors deletes authors with a first name of "joe".
	DeleteAuthors(ctx context.Context) (pgconn.CommandTag, error)

	// DeleteAuthorsByFirstName deletes authors by first name.
	DeleteAuthorsByFirstName(ctx context.Context, firstName string) (pgconn.CommandTag, error)

	// DeleteAuthorsByFullName deletes authors by the full name.
	DeleteAuthorsByFullName(ctx context.Context, params DeleteAuthorsByFullNameParams) (pgconn.CommandTag, error)

	// InsertAuthor inserts an author by name and returns the ID.
	InsertAuthor(ctx context.Context, firstName string, lastName string) (int32, error)

	// InsertAuthorSuffix inserts an author by name and suffix and returns the
	// entire row.
	InsertAuthorSuffix(ctx context.Context, params InsertAuthorSuffixParams) (InsertAuthorSuffixRow, error)

	StringAggFirstName(ctx context.Context, authorID int32) (*string, error)

	ArrayAggFirstName(ctx context.Context, authorID int32) ([]string, error)
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

	LoadType(ctx context.Context, typeName string) (*pgtype.Type, error)
	TypeMap() *pgtype.Map
}

// NewQuerier creates a DBQuerier that implements Querier.
func NewQuerier(ctx context.Context, conn genericConn) (*DBQuerier, error) {
	if err := register(ctx, conn); err != nil {
		return nil, fmt.Errorf("failed to create new querier: %w", err)
	}

	return &DBQuerier{
		conn: conn, 
	}, nil
}

type typeHook func(ctx context.Context, conn genericConn) error

var typeHooks []typeHook

func addHook(hook typeHook) {
	typeHooks = append(typeHooks, hook)
}

func register(ctx context.Context, conn genericConn) error {
	for _, hook := range typeHooks {
		if err := hook(ctx, conn); err != nil {
			return err
		}
	}
	
	return nil
}



const findAuthorByIDSQL = `SELECT * FROM author WHERE author_id = $1;`

type FindAuthorByIDRow struct {
	AuthorID  int32   `json:"author_id"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Suffix    *string `json:"suffix"`
}

// FindAuthorByID implements Querier.FindAuthorByID.
func (q *DBQuerier) FindAuthorByID(ctx context.Context, authorID int32) (FindAuthorByIDRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "FindAuthorByID")
	rows, err := q.conn.Query(ctx, findAuthorByIDSQL, authorID)
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

const findAuthorsSQL = `SELECT * FROM author WHERE first_name = $1;`

type FindAuthorsRow struct {
	AuthorID  int32   `json:"author_id"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Suffix    *string `json:"suffix"`
}

// FindAuthors implements Querier.FindAuthors.
func (q *DBQuerier) FindAuthors(ctx context.Context, firstName string) ([]FindAuthorsRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "FindAuthors")
	rows, err := q.conn.Query(ctx, findAuthorsSQL, firstName)
	if err != nil {
		return nil, fmt.Errorf("query FindAuthors: %w", err)
	}

	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (FindAuthorsRow, error) {
		var item FindAuthorsRow
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

const findAuthorNamesSQL = `SELECT first_name, last_name FROM author ORDER BY author_id = $1;`

type FindAuthorNamesRow struct {
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
}

// FindAuthorNames implements Querier.FindAuthorNames.
func (q *DBQuerier) FindAuthorNames(ctx context.Context, authorID int32) ([]FindAuthorNamesRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "FindAuthorNames")
	rows, err := q.conn.Query(ctx, findAuthorNamesSQL, authorID)
	if err != nil {
		return nil, fmt.Errorf("query FindAuthorNames: %w", err)
	}

	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (FindAuthorNamesRow, error) {
		var item FindAuthorNamesRow
		if err := row.Scan(&item.FirstName, // 'first_name', 'FirstName', '*string', '', '*string'
			&item.LastName, // 'last_name', 'LastName', '*string', '', '*string'
			); err != nil {
			return item, fmt.Errorf("failed to scan: %w", err)
		}
		return item, nil
	})
}

const findFirstNamesSQL = `SELECT first_name FROM author ORDER BY author_id = $1;`

// FindFirstNames implements Querier.FindFirstNames.
func (q *DBQuerier) FindFirstNames(ctx context.Context, authorID int32) ([]*string, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "FindFirstNames")
	rows, err := q.conn.Query(ctx, findFirstNamesSQL, authorID)
	if err != nil {
		return nil, fmt.Errorf("query FindFirstNames: %w", err)
	}

	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (*string, error) {
		var item *string
		if err := row.Scan(&item,
			); err != nil {
			return item, fmt.Errorf("failed to scan: %w", err)
		}
		return item, nil
	})
}

const deleteAuthorsSQL = `DELETE FROM author WHERE first_name = 'joe';`

// DeleteAuthors implements Querier.DeleteAuthors.
func (q *DBQuerier) DeleteAuthors(ctx context.Context) (pgconn.CommandTag, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "DeleteAuthors")
	cmdTag, err := q.conn.Exec(ctx, deleteAuthorsSQL)
	if err != nil {
		return pgconn.CommandTag{}, fmt.Errorf("exec query DeleteAuthors: %w", err)
	}
	return cmdTag, err
}

const deleteAuthorsByFirstNameSQL = `DELETE FROM author WHERE first_name = $1;`

// DeleteAuthorsByFirstName implements Querier.DeleteAuthorsByFirstName.
func (q *DBQuerier) DeleteAuthorsByFirstName(ctx context.Context, firstName string) (pgconn.CommandTag, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "DeleteAuthorsByFirstName")
	cmdTag, err := q.conn.Exec(ctx, deleteAuthorsByFirstNameSQL, firstName)
	if err != nil {
		return pgconn.CommandTag{}, fmt.Errorf("exec query DeleteAuthorsByFirstName: %w", err)
	}
	return cmdTag, err
}

const deleteAuthorsByFullNameSQL = `DELETE
FROM author
WHERE first_name = $1
  AND last_name = $2
  AND suffix = $3;`

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

const insertAuthorSQL = `INSERT INTO author (first_name, last_name)
VALUES ($1, $2)
RETURNING author_id;`

// InsertAuthor implements Querier.InsertAuthor.
func (q *DBQuerier) InsertAuthor(ctx context.Context, firstName string, lastName string) (int32, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "InsertAuthor")
	rows, err := q.conn.Query(ctx, insertAuthorSQL, firstName, lastName)
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

const insertAuthorSuffixSQL = `INSERT INTO author (first_name, last_name, suffix)
VALUES ($1, $2, $3)
RETURNING author_id, first_name, last_name, suffix;`

type InsertAuthorSuffixParams struct {
	FirstName string `json:"FirstName"`
	LastName  string `json:"LastName"`
	Suffix    string `json:"Suffix"`
}

type InsertAuthorSuffixRow struct {
	AuthorID  int32   `json:"author_id"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Suffix    *string `json:"suffix"`
}

// InsertAuthorSuffix implements Querier.InsertAuthorSuffix.
func (q *DBQuerier) InsertAuthorSuffix(ctx context.Context, params InsertAuthorSuffixParams) (InsertAuthorSuffixRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "InsertAuthorSuffix")
	rows, err := q.conn.Query(ctx, insertAuthorSuffixSQL, params.FirstName, params.LastName, params.Suffix)
	if err != nil {
		return InsertAuthorSuffixRow{}, fmt.Errorf("query InsertAuthorSuffix: %w", err)
	}

	return pgx.CollectExactlyOneRow(rows, func(row pgx.CollectableRow) (InsertAuthorSuffixRow, error) {
		var item InsertAuthorSuffixRow
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

const stringAggFirstNameSQL = `SELECT string_agg(first_name, ',') AS names FROM author WHERE author_id = $1;`

// StringAggFirstName implements Querier.StringAggFirstName.
func (q *DBQuerier) StringAggFirstName(ctx context.Context, authorID int32) (*string, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "StringAggFirstName")
	rows, err := q.conn.Query(ctx, stringAggFirstNameSQL, authorID)
	if err != nil {
		return nil, fmt.Errorf("query StringAggFirstName: %w", err)
	}

	return pgx.CollectExactlyOneRow(rows, func(row pgx.CollectableRow) (*string, error) {
		var item *string
		if err := row.Scan(&item,
			); err != nil {
			return item, fmt.Errorf("failed to scan: %w", err)
		}
		return item, nil
	})
}

const arrayAggFirstNameSQL = `SELECT array_agg(first_name) AS names FROM author WHERE author_id = $1;`

// ArrayAggFirstName implements Querier.ArrayAggFirstName.
func (q *DBQuerier) ArrayAggFirstName(ctx context.Context, authorID int32) ([]string, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "ArrayAggFirstName")
	rows, err := q.conn.Query(ctx, arrayAggFirstNameSQL, authorID)
	if err != nil {
		return nil, fmt.Errorf("query ArrayAggFirstName: %w", err)
	}

	return pgx.CollectExactlyOneRow(rows, func(row pgx.CollectableRow) ([]string, error) {
		var item []string
		if err := row.Scan(&item,
			); err != nil {
			return item, fmt.Errorf("failed to scan: %w", err)
		}
		return item, nil
	})
}
