package api

import (
	"strings"

	"github.com/lmnq/grpc-microservices/order/internal/application/core/domain"
	"github.com/lmnq/grpc-microservices/order/internal/ports"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Application struct {
	db      ports.DBPort
	payment ports.PaymentPort
}

func NewApplication(db ports.DBPort, payment ports.PaymentPort) *Application {
	return &Application{
		db:      db,
		payment: payment,
	}
}

func (a Application) PlaceOrder(order domain.Order) (domain.Order, error) {
	err := a.db.Save(&order)
	if err != nil {
		return domain.Order{}, err
	}
	paymentErr := a.payment.Charge(&order)
	if paymentErr != nil {
		st := status.Convert(paymentErr)
		var allErrs []string
		for _, errDetail := range st.Details() {
			switch t := errDetail.(type) {
			case *errdetails.BadRequest:
				for _, violation := range t.GetFieldViolations() {
					allErrs = append(allErrs, violation.Description)
				}
			}
		}
		fieldErr := &errdetails.BadRequest_FieldViolation{
			Field:       "payment",
			Description: strings.Join(allErrs, "\n"),
		}
		badReq := &errdetails.BadRequest{
			// FieldViolations: []*errdetails.BadRequest_FieldViolation{fieldErr},
		}
		badReq.FieldViolations = append(badReq.FieldViolations, fieldErr)
		orderStatus := status.New(codes.InvalidArgument, "order creation failed")
		statusWithDetails, _ := orderStatus.WithDetails(badReq)
		return domain.Order{}, statusWithDetails.Err()
	}

	return order, nil
}
