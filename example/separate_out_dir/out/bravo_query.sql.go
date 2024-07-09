// Code generated by pggen. DO NOT EDIT.

package out

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
)

var _ genericConn = (*pgx.Conn)(nil)
var _ RegisterConn = (*pgx.Conn)(nil)

const bravoSQL = `SELECT 'bravo' as output;`

// Bravo implements Querier.Bravo.
func (q *DBQuerier) Bravo(ctx context.Context) (string, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "Bravo")
	rows, err := q.conn.Query(ctx, bravoSQL)
	if err != nil {
		return "", fmt.Errorf("query Bravo: %w", err)
	}

	return pgx.CollectExactlyOneRow(rows, func(row pgx.CollectableRow) (string, error) {
		var item string
		if err := row.Scan(&item,
			); err != nil {
			return item, fmt.Errorf("failed to scan: %w", err)
		}
		return item, nil
	})
}
