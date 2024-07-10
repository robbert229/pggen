package order

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/robbert229/pggen/internal/pgtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewQuerier_FindOrdersByCustomer(t *testing.T) {
	pool, cleanup := pgtest.NewPostgresSchema(t, []string{"../01_schema.sql", "../02_schema.sql"}, func(config *pgxpool.Config) {
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

	ctx := context.Background()

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
		OrderDate:  pgtype.Timestamptz{Time: time.Now(), Valid: true},
		OrderTotal: pgtype.Numeric{Int: big.NewInt(77), Valid: true},
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
	pool, cleanup := pgtest.NewPostgresSchema(t, []string{"../01_schema.sql", "../02_schema.sql"}, func(config *pgxpool.Config) {
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

	require.NotNil(t, q.FindOrdersByCustomer)
	require.NotNil(t, q.FindProductsInOrder)
	require.NotNil(t, q.InsertOrder)
	require.NotNil(t, q.FindOrdersByPrice)
	require.NotNil(t, q.FindOrdersMRR)
}
