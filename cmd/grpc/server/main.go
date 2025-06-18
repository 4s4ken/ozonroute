package main

import (
	"gitlab.ozon.dev/gleb_evseev/homework/internal/app"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/config"
	hl "gitlab.ozon.dev/gleb_evseev/homework/internal/infra/grpc"
	"log"
	"net"
	"time"

	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/mw"
	pb "gitlab.ozon.dev/gleb_evseev/homework/pkg/pvz"
	"google.golang.org/grpc"
)

const (
	grpcAddress = "localhost:50051"
)

func main() {
	lis, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			mw.LoggingInterceptor,
			mw.ValidateInterceptor,
			mw.RateLimiterInterceptor(
				limiter.New(
					memory.NewStore(),
					limiter.Rate{Period: 10 * time.Second, Limit: 2},
				),
			),
			mw.ErrorMappingInterceptor(),
		),
	)
	js, err := app.NewJsonStorage(config.DBPath, time.Now().UTC)

	if err != nil {
		log.Fatalf("Failed to create json storage: %v", err)
	}

	handler := app.New(js)
	ordersServer := hl.NewOrdersServer(handler)
	pb.RegisterOrdersServiceServer(grpcServer, ordersServer)

	log.Println("gRPC server started on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
