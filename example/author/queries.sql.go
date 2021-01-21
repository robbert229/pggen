// Code generated by pggen. DO NOT EDIT.

package author

import (
	"context"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jschaf/sqld"
)

const findAuthorsSQL = `SELECT first_name, last_name FROM author WHERE first_name = $1`

type Author struct {
	FirstName string
	LastName  string
}

// FindAuthors implements Querier.FindAuthors.
func (q *DBQuerier) FindAuthors(ctx context.Context, firstName string) (auths []Author, mErr error) {
	trace := sqld.ContextClientTrace(ctx)
	traceSendQuery(trace, extractConfig(q.conn), findAuthorsSQL)
	rows, err := q.conn.Query(ctx, findAuthorsSQL, firstName)
	cmdTag := pgconn.CommandTag{}
	if rows != nil {
		cmdTag = rows.CommandTag()
		defer rows.Close()
	}
	traceGotResponse(trace, rows, cmdTag, err)
	defer traceScanResponse(trace, mErr)
	if err != nil {
		return nil, fmt.Errorf("query FindAuthors: %w", err)
	}
	var items []Author
	for rows.Next() {
		var item Author
		if err := rows.Scan(&item.FirstName, &item.LastName); err != nil {
			return nil, fmt.Errorf("scan FindAuthors row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, err
}

// FindAuthorsBatch implements Querier.FindAuthorsBatch.
func (q *DBQuerier) FindAuthorsBatch(ctx context.Context, batch pgx.Batch, firstName string) {
	traceEnqueueQuery(sqld.ContextClientTrace(ctx), findAuthorsSQL)
	batch.Queue(findAuthorsSQL, firstName)
}

// FindAuthorsScan implements Querier.FindAuthorsScan.
func (q *DBQuerier) FindAuthorsScan(ctx context.Context, results pgx.BatchResults) (auths []Author, mErr error) {
	defer traceScanResponse(sqld.ContextClientTrace(ctx), mErr)
	rows, err := results.Query()
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		return nil, err
	}
	var items []Author
	for rows.Next() {
		var item Author
		if err := rows.Scan(&item.FirstName, &item.LastName); err != nil {
			return nil, fmt.Errorf("scan FindAuthors batch row: %w", err)
		}
		items = append(items, item)
	}
	if rows.Err() != nil {
		return nil, err
	}
	return items, err
}

const deleteAuthorsSQL = `DELETE FROM author where first_name = 'joe'`

// DeleteAuthors implements Querier.DeleteAuthors.
func (q *DBQuerier) DeleteAuthors(ctx context.Context) (pgconn.CommandTag, error) {
	trace := sqld.ContextClientTrace(ctx)
	traceSendQuery(trace, extractConfig(q.conn), findAuthorsSQL)
	cmdTag, err := q.conn.Exec(ctx, deleteAuthorsSQL)
	traceGotResponse(trace, nil, cmdTag, err)
	traceScanResponse(trace, err)
	return cmdTag, err
}

// DeleteAuthorsBatch implements Querier.DeleteAuthorsBatch.
func (q *DBQuerier) DeleteAuthorsBatch(ctx context.Context, batch *pgx.Batch) {
	traceEnqueueQuery(sqld.ContextClientTrace(ctx), deleteAuthorsSQL)
	batch.Queue(deleteAuthorsSQL)
}

// DeleteAuthorsScan implements Querier.DeleteAuthorsScan.
func (q *DBQuerier) DeleteAuthorsScan(ctx context.Context, results pgx.BatchResults) (tag pgconn.CommandTag, mErr error) {
	defer traceScanResponse(sqld.ContextClientTrace(ctx), mErr)
	cmdTag, err := results.Exec()
	return cmdTag, err
}
