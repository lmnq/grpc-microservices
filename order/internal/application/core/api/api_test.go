package api

import (
	"context"
	"errors"
	"testing"

	"github.com/lmnq/grpc-microservices/order/internal/application/core/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type mockDB struct {
	mock.Mock
}

func (m *mockDB) Save(ctx context.Context, order *domain.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func (m *mockDB) Get(ctx context.Context, id int64) (domain.Order, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(domain.Order), args.Error(1)
}

type mockPayment struct {
	mock.Mock
}

func (m *mockPayment) Charge(ctx context.Context, order *domain.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func TestPlaceOrder(t *testing.T) {
	db := &mockDB{}
	payment := &mockPayment{}
	application := NewApplication(db, payment)
	order := domain.Order{
		CustomerID: 1,
		OrderItems: []domain.OrderItem{
			{
				ProductCode: "product-1",
				UnitPrice:   100,
				Quantity:    1,
			},
		},
	}
	db.On("Save", mock.Anything, mock.Anything).Return(nil)
	payment.On("Charge", mock.Anything, mock.Anything).Return(nil)
	result, err := application.PlaceOrder(context.Background(), order)
	assert.NoError(t, err)
	assert.Equal(t, order, result)
}

func TestShouldReturnErrorWhenDBPersistanceFails(t *testing.T) {
	db := &mockDB{}
	payment := &mockPayment{}
	application := NewApplication(db, payment)
	order := domain.Order{
		CustomerID: 1,
		OrderItems: []domain.OrderItem{
			{
				ProductCode: "product-1",
				UnitPrice:   100,
				Quantity:    1,
			},
		},
	}
	payment.On("Charge", mock.Anything, mock.Anything).Return(nil)
	db.On("Save", mock.Anything, mock.Anything).Return(errors.New("connection error"))
	_, err := application.PlaceOrder(context.Background(), order)
	assert.EqualError(t, err, "connection error")
}

func TestShouldReturnErrorWhenPaymentFails(t *testing.T) {
	db := &mockDB{}
	payment := &mockPayment{}
	application := NewApplication(db, payment)
	order := domain.Order{
		CustomerID: 1,
		OrderItems: []domain.OrderItem{
			{
				ProductCode: "product-1",
				UnitPrice:   100,
				Quantity:    1,
			},
		},
	}
	db.On("Save", mock.Anything, mock.Anything).Return(nil)
	payment.On("Charge", mock.Anything, mock.Anything).Return(errors.New("insufficient balance"))
	_, err := application.PlaceOrder(context.Background(), order)
	st, _ := status.FromError(err)
	assert.Equal(t, st.Code(), codes.InvalidArgument)
	assert.Equal(t, st.Message(), "order creation failed")
	assert.Equal(t, st.Details()[0].(*errdetails.BadRequest).FieldViolations[0].Description, "insufficient balance")
}
