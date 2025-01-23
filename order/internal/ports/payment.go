package ports

import (
	"context"

	"github.com/lmnq/grpc-microservices/order/internal/application/core/domain"
)

type PaymentPort interface {
	Charge(ctx context.Context, order *domain.Order) error
}
