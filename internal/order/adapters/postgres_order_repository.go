package adapters

import (
	"context"

	domain "github.com/MarinDmitrii/training-store-backend/internal/order/domain"
	"github.com/jmoiron/sqlx"
)

type PostgresOrderRepository struct {
	db *sqlx.DB
}

func NewPostgresOrderRepository(db *sqlx.DB) *PostgresOrderRepository {
	return &PostgresOrderRepository{db: db}
}

func (r *PostgresOrderRepository) SaveOrder(ctx context.Context, domainOrder domain.Order) (int, error) {
	query := `
		INSERT INTO orders (id, payment_key, status, total_price)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO UPDATE 
		SET payment_key = EXCLUDED.payment_key, status = EXCLUDED.status, total_price = EXCLUDED.total_price
		RETURNING id
	`

	var orderId int
	err := r.db.QueryRowContext(ctx, query, domainOrder.Id, domainOrder.Payment_key, domainOrder.Status, domainOrder.Total_price).Scan(&orderId)
	if err != nil {
		return 0, err
	}

	return orderId, nil
}

func (r *PostgresOrderRepository) GetOrders(ctx context.Context) ([]domain.Order, error) {
	query := `SELECT * FROM orders`

	var orders []domain.Order
	err := r.db.SelectContext(ctx, &orders, query)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *PostgresOrderRepository) GetOrderById(ctx context.Context, orderId int) (domain.Order, error) {
	query := "SELECT * FROM orders WHERE id = $1"

	var order domain.Order
	err := r.db.GetContext(ctx, &order, query, orderId)
	if err != nil {
		return domain.Order{}, err
	}

	return order, nil
}

func (r *PostgresOrderRepository) GetOrderByPaymentKey(ctx context.Context, payment_key string) (domain.Order, error) {
	query := `SELECT * FROM orders WHERE payment_key = $1`

	var order domain.Order
	err := r.db.GetContext(ctx, &order, query, payment_key)
	if err != nil {
		return domain.Order{}, err
	}

	return order, nil
}
