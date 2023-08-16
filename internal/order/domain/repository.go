package domain

import (
	"context"
)

type Repository interface {
	SaveOrder(ctx context.Context, order Order) (int, error)
	GetOrders(ctx context.Context) ([]Order, error)
	GetOrderById(ctx context.Context, OrderId int) (Order, error)
	GetOrderByPaymentKey(ctx context.Context, payment_key string) (Order, error)
}
