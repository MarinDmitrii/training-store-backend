package adapters

import (
	"context"
	"os"

	"github.com/MarinDmitrii/training-store-backend/common/log"
	"github.com/MarinDmitrii/training-store-backend/internal/order/usecase"
	"github.com/joho/godotenv"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/checkout/session"
)

type StripeService struct {
	key        string
	successURL string
	cancelURL  string
}

func NewStripeService() *StripeService {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	key := os.Getenv("STRIPE_PRIVATE_KEY")
	successURL := os.Getenv("STRIPE_SUCCESS_URL")
	cancelURL := os.Getenv("STRIPE_CANCEL_URL")

	stripe.Key = key
	stripe.DefaultLeveledLogger = log.Logger

	return &StripeService{
		key:        key,
		successURL: successURL,
		cancelURL:  cancelURL,
	}
}

func (s *StripeService) CreatePayment(ctx context.Context, order usecase.SaveOrder) (*usecase.PaymentResponse, error) {
	lineItems := make([]*stripe.CheckoutSessionLineItemParams, len(order.Products))

	for i, product := range order.Products {
		lineItem := &stripe.CheckoutSessionLineItemParams{
			PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
				Currency: stripe.String(string(stripe.CurrencyUSD)),
				ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
					Name:   stripe.String(product.Name),
					Images: []*string{stripe.String(product.Image)},
				},
				UnitAmount: stripe.Int64(int64(product.Price * 100)),
			},
			Quantity: stripe.Int64(int64(product.Quantity)),
		}

		lineItems[i] = lineItem
	}

	params := &stripe.CheckoutSessionParams{
		LineItems:  lineItems,
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(s.successURL),
		CancelURL:  stripe.String(s.cancelURL),
	}

	response, err := session.New(params)
	if err != nil {
		return nil, err
	}

	return &usecase.PaymentResponse{
		Key: response.ID,
		URL: response.URL,
	}, nil
}
