package main

import (
	"context"
	"log"
	"sync"
	"time"

	usecase "gitlab.ozon.dev/gleb_evseev/homework/pkg/pvz"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	//now := time.Date(2027, 6, 26, 12, 0, 0, 0, time.UTC)
	issuers := map[uint64][]uint64{
		13: {505},
	}
	var wg sync.WaitGroup
	for userId, orderIds := range issuers {
		wg.Add(1)
		go func(uid uint64, oids []uint64) {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			req := &usecase.ProcessOrdersRequest{
				UserId:   uid,
				Action:   usecase.ActionType_ACTION_TYPE_ISSUE,
				OrderIds: oids,
			}

			resp, err := client.ProcessOrders(ctx, req)
			if err != nil {
				log.Printf("ProcessOrders failed for user %d: %v", uid, err)
				return
			}
			log.Printf("ProcessOrders success for user %d: %v", uid, resp)
		}(userId, orderIds)
	}

	wg.Wait()
}
