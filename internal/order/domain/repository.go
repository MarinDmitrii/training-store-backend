package domain

import (
	"context"
)

type Repository interface {
	SaveOrder(ctx context.Context, order Order) (int, error)
	GetOrders(ctx context.Context) ([]Order, error)
	GetOrderByID(ctx context.Context, OrderID int) (Order, error)
	GetOrderByPaymentKey(ctx context.Context, payment_key string) (Order, error)
}
