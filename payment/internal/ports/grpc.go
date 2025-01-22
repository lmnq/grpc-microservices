package ports

import "github.com/lmnq/grpc-microservices/payment/internal/application/domain"

type APIPort interface {
	Charge(payment domain.Payment) (domain.Payment, error)
}
