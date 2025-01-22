package main

import (
	"log"

	"github.com/lmnq/grpc-microservices/payment/config"
	"github.com/lmnq/grpc-microservices/payment/internal/adapters/db"
	"github.com/lmnq/grpc-microservices/payment/internal/adapters/grpc"
	"github.com/lmnq/grpc-microservices/payment/internal/application/api"
)

func main() {
	dbAdapter, err := db.NewAdapter(config.GetDataSourceURL())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	application := api.NewApplication(dbAdapter)
	grpcServer := grpc.NewAdapter(application, config.GetApplicationPort())
	log.Printf("starting grpc server on port %d", config.GetApplicationPort())
	grpcServer.Run()
}
