package usecase

import (
	"context"

	"github.com/MarinDmitrii/training-store-backend/internal/order/domain"
)

type GetOrderByIdUseCase struct {
	orderRepository domain.Repository
}

func NewGetOrderByIdUseCase(orderRepository domain.Repository) *GetOrderByIdUseCase {
	return &GetOrderByIdUseCase{orderRepository: orderRepository}
}

func (uc *GetOrderByIdUseCase) Execute(ctx context.Context, orderId int) (domain.Order, error) {
	return uc.orderRepository.GetOrderByID(ctx, orderId)
}
