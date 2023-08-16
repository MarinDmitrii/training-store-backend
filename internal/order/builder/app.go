package builder

import (
	"context"
	"fmt"

	"github.com/MarinDmitrii/training-store-backend/internal/order/adapters"
	"github.com/MarinDmitrii/training-store-backend/internal/order/usecase"

	"github.com/jmoiron/sqlx"
)

type Application struct {
	SaveOrder      *usecase.SaveOrderUseCase
	GetOrders      *usecase.GetOrdersUseCase
	ConfirmPayment *usecase.ConfirmPaymentUseCase
}

func NewApplication(ctx context.Context) (*Application, func()) {
	cfg, err := NewPostgresConfig()
	if err != nil {
		panic(err)
	}

	postgresConnect := fmt.Sprintf("host= %s port= %d user= %s password= %s dbname= %s sslmode=disabled",
		cfg.Host, cfg.Port, cfg.User, cfg.password, cfg.Database)
	db, err := sqlx.ConnectContext(ctx, "postgres", postgresConnect)
	// db, err := sqlx.Open("postgres", postgresConnect)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	orderRepository := adapters.NewPostgresOrderRepository(db)

	return &Application{
		SaveOrder:      usecase.NewSaveOrderUseCase(orderRepository),
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
