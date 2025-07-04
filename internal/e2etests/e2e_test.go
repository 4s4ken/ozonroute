package e2e

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"

	pkg "gitlab.ozon.dev/gleb_evseev/homework/pkg/pvz"
)

const grpcPort = "50051"

func startGrpcServer(ctx context.Context, dbDSN string) (*exec.Cmd, error) {
	cmd := exec.CommandContext(ctx, "./server.exe", "-db", dbDSN, "-port", grpcPort)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		return nil, err
	}
	return cmd, nil
}

func TestE2E(t *testing.T) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:16",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "testuser",
			"POSTGRES_PASSWORD": "testpass",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp").WithStartupTimeout(60 * time.Second),
	}
	pgContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("failed to start postgres container: %v", err)
	}
	defer func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Logf("failed to terminate pgContainer: %v", err)
		}
	}()

	host, err := pgContainer.Host(ctx)
	if err != nil {
		t.Fatalf("failed to get container host: %v", err)
	}
	port, err := pgContainer.MappedPort(ctx, "5432")
	if err != nil {
		t.Fatalf("failed to get container port: %v", err)
	}

	dsn := fmt.Sprintf("postgres://testuser:testpass@%s:%s/testdb?sslmode=disable", host, port.Port())

	serverCtx, serverCancel := context.WithCancel(ctx)
	defer serverCancel()

	grpcCmd, err := startGrpcServer(serverCtx, dsn)
	if err != nil {
		t.Fatalf("failed to start grpc server: %v", err)
	}
	defer func() {
		if err := grpcCmd.Process.Kill(); err != nil {
			t.Logf("failed to kill grpc server: %v", err)
		}
	}()

	waitForGRPC(t, grpcPort)

	conn, err := grpc.NewClient("localhost:"+grpcPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("failed to dial grpc: %v", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			t.Logf("failed to close gRPC connection: %v", err)
		}
	}()

	client := pkg.NewOrdersServiceClient(conn)

	t.Run("AcceptOrder", func(t *testing.T) {
		req := &pkg.AcceptOrderRequest{
			OrderId:   123,
			UserId:    123,
			ExpiresAt: timestamppb.New(time.Now().Add(24 * time.Hour)),
			Weight:    10,
			Price:     100,
		}
		resp, err := client.AcceptOrder(ctx, req)
		if err != nil {
			t.Fatalf("AcceptOrder failed: %v", err)
		}
		if resp.OrderId != req.OrderId {
			t.Errorf("expected order id %d, got %d", req.OrderId, resp.OrderId)
		}
	})

	t.Run("ProcessOrders_Issue", func(t *testing.T) {
		req := &pkg.ProcessOrdersRequest{
			UserId:   123,
			Action:   pkg.ActionType_ACTION_TYPE_ISSUE,
			OrderIds: []uint64{123},
		}
		resp, err := client.ProcessOrders(ctx, req)
		if err != nil {
			t.Fatalf("ProcessOrders failed: %v", err)
		}
		found := false
		for _, id := range resp.Processed {
			if id == 123 {
				found = true
			}
		}
		if !found {
			t.Errorf("order %d not processed", 123)
		}
	})
}

func waitForGRPC(t *testing.T, port string) {
	for i := 0; i < 30; i++ {
		conn, err := net.DialTimeout("tcp", "localhost:"+port, 1*time.Second)
		if err == nil {
			_ = conn.Close()
			return
		}
		time.Sleep(1 * time.Second)
	}
	t.Fatal("grpc server not responding on port", port)
}
