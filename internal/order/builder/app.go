package builder

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/MarinDmitrii/training-store-backend/internal/order/adapters"
	"github.com/MarinDmitrii/training-store-backend/internal/order/usecase"
	"github.com/joho/godotenv"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Application struct {
	SaveOrder      *usecase.SaveOrderUseCase
	GetOrders      *usecase.GetOrdersUseCase
	ConfirmPayment *usecase.ConfirmPaymentUseCase
}

type PostgresConfig struct {
	host     string
	port     string
	user     string
	password string
	dbname   string
}

func NewPostgresConfig() *PostgresConfig {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")
	return &PostgresConfig{
		host:     host,
		port:     port,
		user:     user,
		password: password,
		dbname:   dbname,
	}
}

func NewApplication(ctx context.Context) (*Application, func()) {
	PostgresConfig := NewPostgresConfig()
	postgresConnect := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		PostgresConfig.host,
		PostgresConfig.port,
		PostgresConfig.user,
		PostgresConfig.password,
		PostgresConfig.dbname,
	)

	db, err := sqlx.ConnectContext(ctx, "postgres", postgresConnect)
	if err != nil {
		panic(err)
	}

	orderRepository := adapters.NewPostgresOrderRepository(db)
	stripeService := adapters.NewStripeService()

	return &Application{
		SaveOrder:      usecase.NewSaveOrderUseCase(orderRepository, stripeService),
		GetOrders:      usecase.NewGetOrdersUseCase(orderRepository),
		ConfirmPayment: usecase.NewConfirmPaymentUseCase(orderRepository),
	}, func() { db.Close() }
}

// func getPaymentService() usecase.PaymentService {
// 	switch os.Getenv("PAYMENT_SERVICE") {
// 	case "stripe":
// 		return adapters.NewStripeService(
// 			os.Getenv("STRIPE_PRIVATE_KEY"),
// 			os.Getenv("STRIPE_TRANSACTION_URL"),
// 			os.Getenv("STRIPE_SUCCESS_URL"),
// 			os.Getenv("STRIPE_CANCEL_URL"),
// 		)
// 	}

// 	log.Logger.Fatal("unknown payment type")
// 	return nil
// }
