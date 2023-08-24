package domain

import "errors"

const (
	PaymentStatusCreated = "created"
	PaymentStatusPaid    = "paid"
)

var NotFoundError = errors.New("Order not found")

type Order struct {
	ID          int
	Payment_key string
	Status      string
	Total_price int
}

func (o *Order) ChangeStatus(status string) {
	o.Status = status
}

func (o *Order) CompletePayment() {
	o.ChangeStatus(PaymentStatusPaid)
}
