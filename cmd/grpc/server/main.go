package main

import (
	"log"
	"net"
	"time"

	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/app"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/mw"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/service/pvz"
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
		),
	)

	ordersService := pvz.NewOrdersService()
	ordersServer := app.NewOrdersServer(ordersService)

	pb.RegisterOrdersServiceServer(grpcServer, ordersServer)

	log.Println("gRPC server started on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
