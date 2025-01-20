package grpc

import (
	"context"

	"github.com/lmnq/grpc-microservices/order/internal/application/core/domain"
	pbOrder "github.com/lmnq/microservices-proto/generated/golang/order"
)

func (a Adapter) Create(ctx context.Context, req *pbOrder.CreateOrderRequest) (*pbOrder.CreateOrderResponse, error) {
	var orderItems []domain.OrderItem
	for _, orderItem := range req.OrderItems {
		orderItems = append(orderItems, domain.OrderItem{
			ProductCode: orderItem.ProductCode,
			UnitPrice:   orderItem.UnitPrice,
			Quantity:    orderItem.Quantity,
		})
	}
	newOrder := domain.NewOrder(req.UserId, orderItems)
	result, err := a.api.PlaceOrder(newOrder)
	if err != nil {
		return nil, err
	}
	return &pbOrder.CreateOrderResponse{
		OrderId: result.ID,
	}, nil
}
