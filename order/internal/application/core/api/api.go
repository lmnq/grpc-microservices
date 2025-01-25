package api

import (
	"context"

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

func (a Application) PlaceOrder(ctx context.Context, order domain.Order) (domain.Order, error) {
	err := a.db.Save(ctx, &order)
	if err != nil {
		return domain.Order{}, err
	}
	paymentErr := a.payment.Charge(ctx, &order)
	if paymentErr != nil {
		// st := status.Convert(paymentErr)
		// var allErrs []string
		// for _, errDetail := range st.Details() {
		// 	switch t := errDetail.(type) {
		// 	case *errdetails.BadRequest:
		// 		for _, violation := range t.GetFieldViolations() {
		// 			allErrs = append(allErrs, violation.Description)
		// 		}
		// 	}
		// }
		// fieldErr := &errdetails.BadRequest_FieldViolation{
		// 	Field:       "payment",
		// 	Description: strings.Join(allErrs, "\n"),
		// }
		st := status.Convert(paymentErr)
		fieldErr := &errdetails.BadRequest_FieldViolation{
			Field:       "payment",
			Description: st.Message(),
		}
		badReq := &errdetails.BadRequest{}
		badReq.FieldViolations = append(badReq.FieldViolations, fieldErr)
		orderStatus := status.New(codes.InvalidArgument, "order creation failed")
		statusWithDetails, _ := orderStatus.WithDetails(badReq)
		return domain.Order{}, statusWithDetails.Err()
	}

	return order, nil
}

func (a Application) GetOrder(ctx context.Context, orderID int64) (domain.Order, error) {
	return a.db.Get(ctx, orderID)
}
