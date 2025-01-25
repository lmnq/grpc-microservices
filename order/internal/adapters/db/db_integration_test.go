package db

import (
	"context"
	"log"
	"testing"

	"github.com/lmnq/grpc-microservices/order/internal/application/core/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type OrderDBTestSuite struct {
	suite.Suite
	ctx            context.Context
	mysqlContainer *MySQLContainer
	dbAdapter      *Adapter
}

func (s *OrderDBTestSuite) SetupSuite() {
	s.ctx = context.Background()
	mysqlContainer, err := CreateMySqlContainer(s.ctx)
	if err != nil {
		log.Fatal(err)
	}
	s.mysqlContainer = mysqlContainer

	dbAdapter, err := NewAdapter(mysqlContainer.DataSourceURL)
	if err != nil {
		log.Fatal(err)
	}
	s.dbAdapter = dbAdapter
}

func (s *OrderDBTestSuite) TearDownSuite() {
	if err := s.mysqlContainer.Terminate(s.ctx); err != nil {
		log.Fatalf("error terminating mysql container: %s", err)
	}
}

func (s *OrderDBTestSuite) TestSaveOrder() {
	t := s.T()
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
	err := s.dbAdapter.Save(s.ctx, &order)
	assert.NoError(t, err)
}

func (s *OrderDBTestSuite) TestGetOrder() {
	t := s.T()
	order := domain.Order{
		CustomerID: 2,
		OrderItems: []domain.OrderItem{
			{
				ProductCode: "product-2",
				UnitPrice:   200,
				Quantity:    2,
			},
		},
	}
	s.dbAdapter.Save(s.ctx, &order)
	savedOrder, err := s.dbAdapter.Get(s.ctx, order.ID)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), savedOrder.CustomerID)
}

func TestOrderDBTestSuite(t *testing.T) {
	suite.Run(t, new(OrderDBTestSuite))
}
