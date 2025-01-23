package ports

import (
	"context"

	"github.com/lmnq/grpc-microservices/order/internal/application/core/domain"
)

type DBPort interface {
	Get(ctx context.Context, id string) (domain.Order, error)
	Save(ctx context.Context, order *domain.Order) error
}
