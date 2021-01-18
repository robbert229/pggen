package author

import (
	"context"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

// QueryHook provides a before and after hook for executing all queries
// generated by sqld.
type QueryHook interface {
	// BeforeQuery runs when a Querier is first invoked.
	BeforeQuery(ctx context.Context, params BeforeHookParams) context.Context

	// AfterQuery runs immediately after the Querier receives the results from
	// executing a query but before the rows are scanned.
	AfterQuery(ctx context.Context, params AfterHookParams)
}

type BeforeHookParams struct {
	// The name of the query as written in the SQL source file.
	QueryName string
	// The SQL query as written in the SQL source file.
	SQL string
	// True if this query is a batched query.
	IsBatch bool
}

type AfterHookParams struct {
	// The name of the query as written in the SQL source file.
	QueryName string
	// The SQL query as written in the SQL source file.
	SQL string
	// True if this query is a batched query.
	IsBatch bool
	// The resulting rows of executing a SELECT statement query. Nil for
	// modification queries like UPDATE, INSERT, and DELETE. You must not retain
	// a reference to rows; It will be invalid after sqld scans the first row.
	Rows pgx.Rows
	// The command tag for the query. Always set if the query succeeded.
	CommandTag pgconn.CommandTag
	// The result query error, if any.
	QueryErr error
}

type nopHook struct{}

var defaultNopHook = nopHook{}

func NewNopHook() QueryHook                                                           { return defaultNopHook }
func (n nopHook) BeforeQuery(ctx context.Context, _ BeforeHookParams) context.Context { return ctx }
func (n nopHook) AfterQuery(context.Context, AfterHookParams)                         {}
