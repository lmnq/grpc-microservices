package ports

import "github.com/lmnq/grpc-microservices/payment/internal/application/domain"

type DBPort interface {
	Save(payment *domain.Payment) error
}
