// Code generated by pggen. DO NOT EDIT.

package out

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
)

var _ genericConn = (*pgx.Conn)(nil)

const alphaSQL = `SELECT 'alpha' as output;`

// Alpha implements Querier.Alpha.
func (q *DBQuerier) Alpha(ctx context.Context) (string, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "Alpha")
	rows, err := q.conn.Query(ctx, alphaSQL)
	if err != nil {
		return "", fmt.Errorf("query Alpha: %w", err)
	}

	return pgx.CollectExactlyOneRow(rows, func(row pgx.CollectableRow) (string, error) {
		var item string
		if err := row.Scan(
			&item,
		); err != nil {
			return item, fmt.Errorf("failed to scan: %w", err)
		}
		return item, nil
	})
}
