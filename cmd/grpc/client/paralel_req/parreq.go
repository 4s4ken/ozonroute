package main

import (
	"context"
	"log"
	"sync"
	"time"

	usecase "gitlab.ozon.dev/gleb_evseev/homework/pkg/pvz"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Fatalf("failed to close connection: %v", err)
		}
	}(conn)

	client := usecase.NewOrdersServiceClient(conn)
	now := time.Date(2027, 6, 26, 12, 0, 0, 0, time.UTC)
	var wg sync.WaitGroup
	for i := 101; i <= 103; i++ {
		wg.Add(1)
		go func(id uint64) {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			_, err := client.AcceptOrder(ctx, &usecase.AcceptOrderRequest{
				OrderId:   id,
				UserId:    13,
				ExpiresAt: timestamppb.New(now),
				Package:   nil,
				Weight:    10,
				Price:     10,
			})
			if err != nil {
				log.Printf("AcceptOrder %d failed: %v", id, err)
			} else {
				log.Printf("AcceptOrder %d success", id)
			}
		}(uint64(i))
	}

	wg.Wait()
}
