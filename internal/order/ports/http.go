package ports

import (
	"net/http"

	"github.com/MarinDmitrii/training-store-backend/common/log"
	"github.com/MarinDmitrii/training-store-backend/internal/order/builder"
	"github.com/MarinDmitrii/training-store-backend/internal/order/domain"
	"github.com/MarinDmitrii/training-store-backend/internal/order/usecase"
	"github.com/labstack/echo/v4"
)

type HttpOrderHandler struct {
	app *builder.Application
}

func NewHttpOrderHandler(app *builder.Application) HttpOrderHandler {
	return HttpOrderHandler{app: app}
}

func (h HttpOrderHandler) SaveOrder(ctx echo.Context) error {
	request := &PostOrder{}
	if err := ctx.Bind(request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	saveOrder := &usecase.SaveOrder{
		Products: make([]usecase.Product, len(request.Products)),
	}

	for i, p := range request.Products {
		saveOrder.Products[i] = usecase.Product{
			ID:       p.ID,
			Image:    p.Image,
			Name:     p.Name,
			Price:    p.Price,
			Quantity: p.Quantity,
		}
	}

	paymentResponse, err := h.app.SaveOrder.Execute(ctx.Request().Context(), saveOrder)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, struct {
		Key string `json:"payment_key"`
		URL string `json:"payment_URL"`
	}{
		Key: paymentResponse.Key,
		URL: paymentResponse.URL,
	})
}

func (h HttpOrderHandler) GetOrders(ctx echo.Context) error {
	orders, err := h.app.GetOrders.Execute(ctx.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	response := make([]*Order, 0, len(orders))
	for _, db := range orders {
		ro := h.mapToResponse(db)
		response = append(response, &ro)
	}

	return ctx.JSON(http.StatusOK, response)
}

func (h HttpOrderHandler) ProcessStripeEvent(ctx echo.Context) error {
	log.HttpRequest(ctx.Request())

	event := StripeEvent{}
	if err := ctx.Bind(&event); err != nil {
		log.Errorf("stripe event error: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	switch event.Type {
	case "checkout.session.completed":
		if err := h.app.ConfirmPayment.Execute(ctx.Request().Context(), event.Data.Object.ID); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return ctx.NoContent(http.StatusOK)
}

func (h HttpOrderHandler) mapToResponse(order domain.Order) Order {
	return Order{
		ID:         order.ID,
		PaymentKey: order.Payment_key,
		Status:     OrderStatus(order.Status),
		TotalPrice: float32(order.Total_price) / 100,
	}
}

func CustomRegisterHandlers(router EchoRouter, si ServerInterface) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.GET("/orders", wrapper.GetOrders)
	router.POST("/orders", wrapper.CreateOrder)
	router.POST("/stripe/webhook", wrapper.ProcessStripeEvent)
}
