package e2e

import (
	"context"
	"path/filepath"
	"testing"

	orderpb "github.com/lmnq/microservices-proto/generated/golang/order"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go/modules/compose"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type OrderE2ETestSuite struct {
	suite.Suite
	ctx             context.Context
	compose         compose.ComposeStack
	orderClient     orderpb.OrderClient
	orderClientConn *grpc.ClientConn
}

func (s *OrderE2ETestSuite) SetupSuite() {
	var err error
	s.ctx = context.Background()
	composePath := filepath.Join("resources", "docker-compose.yml")
	// composeFile, err := os.ReadFile(composePath)
	// s.Require().NoError(err)
	s.compose, err = compose.NewDockerCompose(composePath)
	s.Require().NoError(err, "error creating compose")
	err = s.compose.Up(s.ctx, compose.Wait(true))
	s.Require().NoError(err, "error starting docker compose")

	conn, err := grpc.NewClient("localhost:8080",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	s.Require().NoError(err)
	s.orderClientConn = conn

	s.orderClient = orderpb.NewOrderClient(conn)
}

func (s *OrderE2ETestSuite) TearDownSuite() {
	if s.orderClientConn != nil {
		s.Require().NoError(s.orderClientConn.Close(), "error closing order client connection")
	}

	if s.compose != nil {
		err := s.compose.Down(s.ctx,
			compose.RemoveOrphans(true),
			compose.RemoveImagesLocal,
		)
		s.Require().NoError(err, "error stopping docker compose")
	}
}

func (s *OrderE2ETestSuite) TestOrderFlow() {
	createOrderRes, errCreate := s.orderClient.Create(s.ctx, &orderpb.CreateOrderRequest{
		UserId: 21,
		OrderItems: []*orderpb.OrderItem{
			{
				ProductCode: "product-21",
				UnitPrice:   2100,
				Quantity:    21,
			},
		},
	})
	s.Require().NoError(errCreate)

	getOrderRes, errGet := s.orderClient.Get(s.ctx, &orderpb.GetOrderRequest{
		OrderId: createOrderRes.OrderId,
	})
	s.Require().NoError(errGet)
	s.Require().Equal(int64(21), getOrderRes.UserId)

	orderItem := getOrderRes.OrderItems[0]
	s.Require().Equal("product-21", orderItem.ProductCode)
	s.Require().Equal(float32(2100), orderItem.UnitPrice)
	s.Require().Equal(int32(21), orderItem.Quantity)
}

func TestOrderE2ETestSuite(t *testing.T) {
	suite.Run(t, new(OrderE2ETestSuite))
}
