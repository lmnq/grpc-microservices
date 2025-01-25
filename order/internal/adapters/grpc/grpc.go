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
	result, err := a.api.PlaceOrder(ctx, newOrder)
	if err != nil {
		return nil, err
	}
	return &pbOrder.CreateOrderResponse{
		OrderId: result.ID,
	}, nil
}

func (a Adapter) Get(ctx context.Context, request *pbOrder.GetOrderRequest) (*pbOrder.GetOrderResponse, error) {
	result, err := a.api.GetOrder(ctx, request.OrderId)
	var orderItems []*pbOrder.OrderItem
	for _, orderItem := range result.OrderItems {
		orderItems = append(orderItems, &pbOrder.OrderItem{
			ProductCode: orderItem.ProductCode,
			UnitPrice:   orderItem.UnitPrice,
			Quantity:    orderItem.Quantity,
		})
	}
	if err != nil {
		return nil, err
	}
	return &pbOrder.GetOrderResponse{UserId: result.CustomerID, OrderItems: orderItems}, nil
}
