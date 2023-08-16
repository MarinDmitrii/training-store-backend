package usecase

import (
	"context"

	domain "github.com/MarinDmitrii/training-store-backend/internal/order/domain"
)

type SaveOrder struct {
	Products []struct {
		Id       int
		Image    string
		Name     string
		Price    float32
		Quantity int
	}
}

type SaveOrderUseCase struct {
	orderRepository domain.Repository
}

func NewSaveOrderUseCase(
	orderRepository domain.Repository,
) *SaveOrderUseCase {
	return &SaveOrderUseCase{
		orderRepository: orderRepository,
	}
}

func (uc *SaveOrderUseCase) Execute(ctx context.Context, saveOrder *SaveOrder) (int, error) {
	var totalPrice int
	for _, product := range saveOrder.Products {
		totalPrice += int(product.Price*100) * product.Quantity
	}

	order := domain.Order{
		Status:      domain.PaymentStatusCreated,
		Total_price: totalPrice,
	}

	orderId, err := uc.orderRepository.SaveOrder(ctx, order)
	if err != nil {
		return 0, err
	}

	return orderId, nil
}
