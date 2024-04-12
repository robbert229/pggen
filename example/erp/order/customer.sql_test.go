package order

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/robbert229/pggen/internal/pgtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewQuerier_FindOrdersByCustomer(t *testing.T) {
	conn, cleanup := pgtest.NewPostgresSchema(t, []string{"../01_schema.sql", "../02_schema.sql"})
	defer cleanup()
	ctx := context.Background()

	q, err := NewQuerier(context.Background(), conn)
	require.NoError(t, err)
	cust1, err := q.InsertCustomer(ctx, InsertCustomerParams{
		FirstName: "foo_first",
		LastName:  "foo_last",
		Email:     "foo_email",
	})
	if err != nil {
		t.Error(err)
		return
	}
	order1, err := q.InsertOrder(ctx, InsertOrderParams{
		OrderDate:  pgtype.Timestamptz{Time: time.Now()},
		OrderTotal: pgtype.Numeric{Int: big.NewInt(77)},
		CustID:     cust1.CustomerID,
	})
	if err != nil {
		t.Error(err)
		return
	}

	t.Run("FindOrdersByCustomer", func(t *testing.T) {
		orders, err := q.FindOrdersByCustomer(context.Background(), cust1.CustomerID)
		require.NoError(t, err)
		assert.Equal(t, []FindOrdersByCustomerRow{
			{
				OrderID:    order1.OrderID,
				OrderDate:  order1.OrderDate,
				OrderTotal: order1.OrderTotal,
				CustomerID: order1.CustomerID,
			},
		}, orders)
	})
}

func TestNewQuerier_QuerierMatchesDBQuerier(t *testing.T) {
	conn, cleanup := pgtest.NewPostgresSchema(t, []string{"../01_schema.sql", "../02_schema.sql"})
	defer cleanup()

	q, err := NewQuerier(context.Background(), conn)
	require.NoError(t, err)

	require.NotNil(t, q.FindOrdersByCustomer)
	require.NotNil(t, q.FindProductsInOrder)
	require.NotNil(t, q.InsertOrder)
	require.NotNil(t, q.FindOrdersByPrice)
	require.NotNil(t, q.FindOrdersMRR)
}
