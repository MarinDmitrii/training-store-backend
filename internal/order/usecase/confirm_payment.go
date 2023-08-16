package usecase

import (
	"context"

	domain "github.com/MarinDmitrii/training-store-backend/internal/order/domain"
)

type ConfirmPaymentUseCase struct {
	orderRepository domain.Repository
}

func NewConfirmPaymentUseCase(
	orderRepository domain.Repository,
) *ConfirmPaymentUseCase {
	return &ConfirmPaymentUseCase{
		orderRepository: orderRepository,
	}
}

func (uc *ConfirmPaymentUseCase) Execute(ctx context.Context, payment_key string) error {
	order, err := uc.orderRepository.GetOrderByPaymentKey(ctx, payment_key)
	if err != nil {
		return err
	}

	order.CompletePayment()

	return nil
}
