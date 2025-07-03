package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/swaggo/http-swagger"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/app"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/config"
	hl "gitlab.ozon.dev/gleb_evseev/homework/internal/infra/grpc"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/storage/db"
	usecase "gitlab.ozon.dev/gleb_evseev/homework/pkg/pvz"
	"gitlab.ozon.dev/gleb_evseev/homework/pkg/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gitlab.ozon.dev/gleb_evseev/homework/internal/mw"
)

const (
	grpcAddress    = "localhost:50051"
	swaggerAddress = "localhost:7002"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	database := storage.NewPgxPool(ctx, config.DBConnString)
	dbStorage := db.New(database)
	defer func() {
		if err := dbStorage.Close(); err != nil {
			log.Printf("failed to close storage: %v", err)
		}
	}()

	wp := app.NewWorkerPool(ctx, 10)
	defer wp.Shutdown()

	handler := app.NewService(dbStorage, func() time.Time { return time.Now().UTC() }, wp)

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

	usecase.RegisterOrdersServiceServer(grpcServer, hl.NewOrdersServer(handler))
	go func() {
		log.Println("gRPC server started on", grpcAddress)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve gRPC: %v", err)
		}
	}()

	mux := chi.NewMux()
	gw := runtime.NewServeMux()
	err = usecase.RegisterOrdersServiceHandlerFromEndpoint(
		ctx,
		gw,
		grpcAddress,
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	)
	if err != nil {
		log.Fatalf("failed to register gateway: %v", err)
	}

	mux.Mount("/", gw)
	mux.Get("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		data, err := os.ReadFile("./pkg/pvz/pvz.swagger.json")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(data)
	})
	mux.Handle("/swagger/*", httpSwagger.Handler(httpSwagger.URL("/swagger.json")))

	httpServer := &http.Server{
		Addr:    swaggerAddress,
		Handler: mux,
	}

	go func() {
		log.Println("HTTP Gateway + Swagger started on", swaggerAddress)
		if err := http.ListenAndServe(swaggerAddress, mux); err != nil {
			log.Fatalf("failed to serve HTTP: %v", err)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	<-sigCh
	log.Println("Graceful shutdown")

	cancel()
	wp.Shutdown()

	grpcServer.GracefulStop()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer shutdownCancel()
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("Ошибка при остановке HTTP сервера: %v", err)
	}

	if err := dbStorage.Close(); err != nil {
		log.Printf("Ошибка при закрытии БД: %v", err)
	}
	log.Println("Graceful shutdown done")

}
