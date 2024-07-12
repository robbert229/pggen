// Code generated by pggen. DO NOT EDIT.

package casting

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
	// FindBooks finds all books.
	FindBooks(ctx context.Context) ([]FindBooksRow, error)

	// InsertAuthor inserts an author by name and returns the ID.
	InsertAuthor(ctx context.Context, firstName string, lastName string) (int32, error)

	// InsertBook inserts a book.
	InsertBook(ctx context.Context, title string) (int32, error)

	// AssignAuthor assigns an author to a book.
	AssignAuthor(ctx context.Context, authorID int32, bookID int32) (pgconn.CommandTag, error)
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



// Author represents the Postgres composite type "author".
type Author struct {
	AuthorID  *int32  `json:"author_id"`
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Suffix    *string `json:"suffix"`
}

// Book represents the Postgres composite type "book".
type Book struct {
	BookID   *int32  `json:"book_id"`
	Title    *string `json:"title"`
	AuthorID *int32  `json:"author_id"`
}




	// codec_newAuthor is a codec for the composite type of the same name
	func codec_newAuthor(conn RegisterConn) (pgtype.Codec, error) {
		
		    field0, ok := conn.TypeMap().TypeForName("int4")
			if !ok {
				return nil, fmt.Errorf("type not found: int4")
			}
		
		    field1, ok := conn.TypeMap().TypeForName("text")
			if !ok {
				return nil, fmt.Errorf("type not found: text")
			}
		
		    field2, ok := conn.TypeMap().TypeForName("text")
			if !ok {
				return nil, fmt.Errorf("type not found: text")
			}
		
		    field3, ok := conn.TypeMap().TypeForName("text")
			if !ok {
				return nil, fmt.Errorf("type not found: text")
			}
		
		
		return &pgtype.CompositeCodec{
			Fields: []pgtype.CompositeCodecField{
				
					{
						Name: "author_id",
						Type: field0,
					},
				
					{
						Name: "first_name",
						Type: field1,
					},
				
					{
						Name: "last_name",
						Type: field2,
					},
				
					{
						Name: "suffix",
						Type: field3,
					},
				
			},
		}, nil
	}

	func register_newAuthor(
		ctx context.Context,
		conn RegisterConn,
	) error {
		t, err := conn.LoadType(
			ctx,
			"\"author\"",
		)
		if err != nil {
			return fmt.Errorf("newAuthor failed to load type: %w", err)
		}
		
		conn.TypeMap().RegisterType(t)

		return nil
	}

	func init(){
		addHook(register_newAuthor) 
	}
	


	// codec_newBook is a codec for the composite type of the same name
	func codec_newBook(conn RegisterConn) (pgtype.Codec, error) {
		
		    field0, ok := conn.TypeMap().TypeForName("int4")
			if !ok {
				return nil, fmt.Errorf("type not found: int4")
			}
		
		    field1, ok := conn.TypeMap().TypeForName("text")
			if !ok {
				return nil, fmt.Errorf("type not found: text")
			}
		
		    field2, ok := conn.TypeMap().TypeForName("int4")
			if !ok {
				return nil, fmt.Errorf("type not found: int4")
			}
		
		
		return &pgtype.CompositeCodec{
			Fields: []pgtype.CompositeCodecField{
				
					{
						Name: "book_id",
						Type: field0,
					},
				
					{
						Name: "title",
						Type: field1,
					},
				
					{
						Name: "author_id",
						Type: field2,
					},
				
			},
		}, nil
	}

	func register_newBook(
		ctx context.Context,
		conn RegisterConn,
	) error {
		t, err := conn.LoadType(
			ctx,
			"\"book\"",
		)
		if err != nil {
			return fmt.Errorf("newBook failed to load type: %w", err)
		}
		
		conn.TypeMap().RegisterType(t)

		return nil
	}

	func init(){
		addHook(register_newBook) 
	}
	

const findBooksSQL = `SELECT a::author, b::book
FROM book b LEFT JOIN author a ON (b.author_id = a.author_id);`

type FindBooksRow struct {
	A *Author `json:"a"`
	B *Book   `json:"b"`
}

// FindBooks implements Querier.FindBooks.
func (q *DBQuerier) FindBooks(ctx context.Context) ([]FindBooksRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "FindBooks")
	rows, err := q.conn.Query(ctx, findBooksSQL)
	if err != nil {
		return nil, fmt.Errorf("query FindBooks: %w", err)
	}

	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (FindBooksRow, error) {
  var item FindBooksRow
		if err := row.Scan(&item.A, // 'a', 'A', '*Author', '', '*Author'
			&item.B, // 'b', 'B', '*Book', '', '*Book'
			); err != nil {
			return item, fmt.Errorf("failed to scan: %w", err)
		}
		return item, nil
	})
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

const insertBookSQL = `INSERT INTO book (title)
VALUES ($1)
RETURNING book_id;`

// InsertBook implements Querier.InsertBook.
func (q *DBQuerier) InsertBook(ctx context.Context, title string) (int32, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "InsertBook")
	rows, err := q.conn.Query(ctx, insertBookSQL, title)
	if err != nil {
		return 0, fmt.Errorf("query InsertBook: %w", err)
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

const assignAuthorSQL = `UPDATE book
SET author_id = $1
WHERE book_id = $2;`

// AssignAuthor implements Querier.AssignAuthor.
func (q *DBQuerier) AssignAuthor(ctx context.Context, authorID int32, bookID int32) (pgconn.CommandTag, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "AssignAuthor")
	cmdTag, err := q.conn.Exec(ctx, assignAuthorSQL, authorID, bookID)
	if err != nil {
		return pgconn.CommandTag{}, fmt.Errorf("exec query AssignAuthor: %w", err)
	}
	return cmdTag, err
}
