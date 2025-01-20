package main

import (
	"log"

	"github.com/lmnq/grpc-microservices/order/config"
	"github.com/lmnq/grpc-microservices/order/internal/adapters/db"
	"github.com/lmnq/grpc-microservices/order/internal/adapters/grpc"
	"github.com/lmnq/grpc-microservices/order/internal/application/core/api"
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
