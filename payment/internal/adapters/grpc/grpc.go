package grpc

import (
	"context"
	"fmt"

	"github.com/lmnq/grpc-microservices/payment/internal/application/domain"
	pbPayment "github.com/lmnq/microservices-proto/generated/golang/payment"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a Adapter) Create(ctx context.Context, req *pbPayment.CreatePaymentRequest) (*pbPayment.CreatePaymentResponse, error) {
	newPayment := domain.NewPayment(req.UserId, req.OrderId, req.TotalPrice)
	result, err := a.api.Charge(newPayment)
	if err != nil {
		return nil, status.New(codes.Internal, fmt.Sprintf("failed to charge: %v", err)).Err()
	}
	return &pbPayment.CreatePaymentResponse{
		PaymentId: result.ID,
	}, nil
}
