package usecase

import (
	"context"

	domain "github.com/MarinDmitrii/training-store-backend/internal/order/domain"
)

type Product struct {
	ID       int
	Image    string
	Name     string
	Price    float32
	Quantity int
}

type SaveOrder struct {
	Products []Product
}

type PaymentResponse struct {
	Key string
	URL string
}

type PaymentService interface {
	CreatePayment(context.Context, SaveOrder) (*PaymentResponse, error)
}

type SaveOrderUseCase struct {
	orderRepository domain.Repository
	stripeService   PaymentService
}

func NewSaveOrderUseCase(
	orderRepository domain.Repository,
	stripeService PaymentService,
) *SaveOrderUseCase {
	return &SaveOrderUseCase{
		orderRepository: orderRepository,
		stripeService:   stripeService,
	}
}

func (uc *SaveOrderUseCase) Execute(ctx context.Context, saveOrder *SaveOrder) (*PaymentResponse, error) {
	var totalPrice int
	for _, product := range saveOrder.Products {
		totalPrice += int(product.Price*100) * product.Quantity
	}

	order := domain.Order{
		Status:      domain.PaymentStatusCreated,
		Total_price: totalPrice,
	}

	orderID, err := uc.orderRepository.SaveOrder(ctx, order)
	if err != nil {
		return nil, err
	}
	order.ID = orderID

	paymentResponse, err := uc.stripeService.CreatePayment(ctx, *saveOrder)
	if err != nil {
		return nil, err
	}
	order.Payment_key = paymentResponse.Key

	orderID, err = uc.orderRepository.SaveOrder(ctx, order)
	if err != nil {
		return nil, err
	}

	return &PaymentResponse{
		Key: paymentResponse.Key,
		URL: paymentResponse.URL,
	}, nil
}
