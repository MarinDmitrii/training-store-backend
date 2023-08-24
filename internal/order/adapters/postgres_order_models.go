package adapters

import (
	domain "github.com/MarinDmitrii/training-store-backend/internal/order/domain"
)

type OrderModel struct {
	ID          int    `db:"id"`
	Payment_key string `db:"payment_key"`
	Status      string `db:"status"`
	Total_price int    `db:"total_price"`
}

func NewOrderModel(order domain.Order) (OrderModel, error) {
	return OrderModel{
		ID:          order.ID,
		Payment_key: order.Payment_key,
		Status:      order.Status,
		Total_price: order.Total_price,
	}, nil
}

func (model *OrderModel) mapToDomain() domain.Order {
	return domain.Order{
		ID:          model.ID,
		Payment_key: model.Payment_key,
		Status:      model.Status,
		Total_price: model.Total_price,
	}
}
