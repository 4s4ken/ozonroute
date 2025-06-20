package main

import (
	"context"
	"gitlab.ozon.dev/gleb_evseev/homework/pkg/pvz"
	"google.golang.org/grpc/grpclog"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	grpcAddress = "localhost:50051"
)

func main() {
	conn, err := grpc.NewClient(grpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer func() {
		if err != nil {
			if err := conn.Close(); err != nil {
				grpclog.Errorf("Failed to close conn to %s: %v", grpcAddress, err)
			}
			return
		}
	}()

	client := usecase.NewOrdersServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	packageType := usecase.PackageType_PACKAGE_TYPE_BOX
	acceptResp, err := client.AcceptOrder(ctx, &usecase.AcceptOrderRequest{
		OrderId:   21,
		UserId:    2,
		ExpiresAt: timestamppb.New(time.Date(2027, 6, 5, 12, 0, 0, 0, time.UTC)),
		Weight:    9,
		Price:     100,
		Package:   &packageType,
	})
	if err != nil {
		log.Fatalf("AcceptOrder failed: %v", err)
	}
	log.Printf("AcceptOrder response: %v", acceptResp)

	listOrdersResp, err := client.ListOrders(ctx, &usecase.ListOrdersRequest{
		UserId: 2,
		InPvz:  true,
		Pagination: &usecase.Pagination{
			Page:        0,
			CountOnPage: 10,
		},
	})
	if err != nil {
		log.Fatalf("ListOrders failed: %v", err)
	}
	log.Printf("ListOrders response: %v", listOrdersResp)

	listReturnsResp, err := client.ListReturns(ctx, &usecase.ListReturnsRequest{
		Pagination: &usecase.Pagination{
			Page:        0,
			CountOnPage: 10,
		},
	})
	if err != nil {
		log.Fatalf("ListReturns failed: %v", err)
	}
	log.Printf("ListReturns response: %v", listReturnsResp)

	getHistoryResp, err := client.GetHistory(ctx, &usecase.GetHistoryRequest{})
	if err != nil {
		log.Fatalf("GetHistory failed: %v", err)
	}
	log.Printf("GetHistory response: %v", getHistoryResp)
}
