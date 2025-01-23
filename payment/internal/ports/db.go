package ports

import (
	"context"

	"github.com/lmnq/grpc-microservices/payment/internal/application/domain"
)

type DBPort interface {
	Save(ctx context.Context, payment *domain.Payment) error
}
