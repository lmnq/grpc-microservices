package grpc

import (
	"fmt"
	"log"
	"net"

	"github.com/lmnq/grpc-microservices/payment/config"
	"github.com/lmnq/grpc-microservices/payment/internal/ports"
	pbPayment "github.com/lmnq/microservices-proto/generated/golang/payment"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Adapter struct {
	api  ports.APIPort
	port int
	pbPayment.UnimplementedPaymentServer
}

func NewAdapter(api ports.APIPort, port int) *Adapter {
	return &Adapter{
		api:  api,
		port: port,
	}
}

func (a Adapter) Run() {
	var err error

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		log.Fatalf("failed to listen on port %d: %v", a.port, err)
	}

	grpcServer := grpc.NewServer()
	pbPayment.RegisterPaymentServer(grpcServer, a)
	// reflection.Register(grpcServer)
	if config.GetEnv() == "dev" {
		reflection.Register(grpcServer)
	}

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve grpc: %v", err)
	}
}
