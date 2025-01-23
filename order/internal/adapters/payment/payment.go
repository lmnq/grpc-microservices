package payment

import (
	"context"

	"github.com/lmnq/grpc-microservices/order/internal/application/core/domain"
	paymentpb "github.com/lmnq/microservices-proto/generated/golang/payment"
)

func (a *Adapter) Charge(ctx context.Context, order *domain.Order) error {
	_, err := a.payment.Create(context.Background(),
		&paymentpb.CreatePaymentRequest{
			UserId:     order.CustomerID,
			OrderId:    order.ID,
			TotalPrice: order.TotalPrice(),
		},
	)
	return err
}
