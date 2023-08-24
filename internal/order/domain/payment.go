package domain

const (
	StatusNew        = "new"
	StatusInProgress = "in_progress"
	StatusPaid       = "paid"
	StatusCancelled  = "cancelled"
)

const (
	TypeIncoming = "incoming"
	TypeOutgoing = "outgoing"
)

type Payment struct {
	ID             string
	OrderID        string
	Status         string
	TransactionKey string
	Link           string
	Price          uint
	Type           string
}

func NewPayment(orderID string, price uint, paymentType string) Payment {
	return Payment{
		OrderID: orderID,
		Status:  StatusNew,
		Price:   price,
		Type:    paymentType,
	}
}

func (p *Payment) ChangeStatus(status string) {
	p.Status = status
}

func (p *Payment) Start(transactionKey, link string) {
	p.Status = StatusInProgress
	p.TransactionKey = transactionKey
	p.Link = link
}
