package ports

import (
	"context"

	"github.com/lmnq/grpc-microservices/payment/internal/application/domain"
)

type APIPort interface {
	Charge(ctx context.Context, payment domain.Payment) (domain.Payment, error)
}
