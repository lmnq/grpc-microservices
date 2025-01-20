package ports

import "github.com/lmnq/grpc-microservices/order/internal/application/core/domain"

type PaymentPort interface {
	Charge(*domain.Order) error
}
