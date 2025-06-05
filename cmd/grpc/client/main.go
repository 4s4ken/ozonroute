package main

import (
	"context"
	"log"
	"time"

	pvz "gitlab.ozon.dev/gleb_evseev/homework/pkg/pvz"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	grpcAddress = "localhost:50051"
)

func main() {
	conn, err := grpc.Dial(grpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	client := pvz.NewOrdersServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	packageType := pvz.PackageType_PACKAGE_TYPE_BOX
	acceptResp, err := client.AcceptOrder(ctx, &pvz.AcceptOrderRequest{
		OrderId:   2,
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

	listOrdersResp, err := client.ListOrders(ctx, &pvz.ListOrdersRequest{
		UserId: 1,
		InPvz:  true,
		Pagination: &pvz.Pagination{
			Page:        1,
			CountOnPage: 10,
		},
	})
	if err != nil {
		log.Fatalf("ListOrders failed: %v", err)
	}
	log.Printf("ListOrders response: %v", listOrdersResp)

	listReturnsResp, err := client.ListReturns(ctx, &pvz.ListReturnsRequest{
		Pagination: &pvz.Pagination{
			Page:        1,
			CountOnPage: 10,
		},
	})
	if err != nil {
		log.Fatalf("ListReturns failed: %v", err)
	}
	log.Printf("ListReturns response: %v", listReturnsResp)

	getHistoryResp, err := client.GetHistory(ctx, &pvz.GetHistoryRequest{
		Pagination: &pvz.Pagination{
			Page:        1,
			CountOnPage: 10,
		},
	})
	if err != nil {
		log.Fatalf("GetHistory failed: %v", err)
	}
	log.Printf("GetHistory response: %v", getHistoryResp)
}
