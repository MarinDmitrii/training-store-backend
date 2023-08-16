package usecase

import (
	"context"

	"github.com/MarinDmitrii/training-store-backend/internal/order/domain"
)

type GetOrdersUseCase struct {
	orderRepository domain.Repository
}

func NewGetOrdersUseCase(orderRepository domain.Repository) *GetOrdersUseCase {
	return &GetOrdersUseCase{orderRepository: orderRepository}
}

func (uc *GetOrdersUseCase) Execute(ctx context.Context) ([]domain.Order, error) {
	return uc.orderRepository.GetOrders(ctx)
}
