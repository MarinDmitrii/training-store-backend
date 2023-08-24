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
	db.MustExec(`
	CREATE TABLE IF NOT EXISTS "product" (
		"id" serial NOT NULL,
		"name" varchar(255) NOT NULL,
		"price" money NOT NULL,
		"quantity" integer NOT NULL,
		CONSTRAINT "product_pk" PRIMARY KEY ("id")
	) WITH (
	  OIDS=FALSE
	);
		
	CREATE TABLE IF NOT EXISTS "orders" (
		"id" serial NOT NULL,
		"payment_key" varchar(700) NOT NULL,
		"status" varchar(255) NOT NULL,
		"total_price" integer NOT NULL,
		CONSTRAINT "order_pk" PRIMARY KEY ("id")
	) WITH (
	  OIDS=FALSE
	);
		
	CREATE TABLE IF NOT EXISTS "order_product" (
		"id" serial NOT NULL,
		"order_id" integer NOT NULL,
		"product_id" integer NOT NULL,
		"quantity" integer NOT NULL,
		CONSTRAINT "order_product_pk" PRIMARY KEY ("id"),
		CONSTRAINT "order_product_fk0" FOREIGN KEY ("order_id") REFERENCES "orders"("id"),
		CONSTRAINT "order_product_fk1" FOREIGN KEY ("product_id") REFERENCES "product"("id")
	) WITH (
	  OIDS=FALSE
	);	
	`)

	return &PostgresOrderRepository{db: db}
}

func (r *PostgresOrderRepository) SaveOrder(ctx context.Context, domainOrder domain.Order) (int, error) {
	if domainOrder.ID == 0 {
		err := r.db.QueryRowContext(ctx, "SELECT nextval('orders_id_seq')").Scan(&domainOrder.ID)
		if err != nil {
			return 0, err
		}
	}

	query := `
		INSERT INTO orders (id, payment_key, status, total_price)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO UPDATE
		SET payment_key = EXCLUDED.payment_key, status = EXCLUDED.status, total_price = EXCLUDED.total_price
		RETURNING id
	`

	var orderID int
	err := r.db.QueryRowContext(ctx, query, domainOrder.ID, domainOrder.Payment_key, domainOrder.Status, domainOrder.Total_price).Scan(&orderID)
	if err != nil {
		return 0, err
	}

	return orderID, nil
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

func (r *PostgresOrderRepository) GetOrderByID(ctx context.Context, orderID int) (domain.Order, error) {
	query := "SELECT * FROM orders WHERE id = $1"

	var order domain.Order
	err := r.db.GetContext(ctx, &order, query, orderID)
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
