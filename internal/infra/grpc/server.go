package grpc

import (
	"log"
	"net"

	desc "gitlab.ozon.dev/gleb_evseev/homework/pkg/pvz"
	"google.golang.org/grpc"
)

func RunGRPCServer(address string, handler desc.OrdersServiceServer) {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	desc.RegisterOrdersServiceServer(s, handler)

	log.Printf("Listening gRPC on %v", address)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
