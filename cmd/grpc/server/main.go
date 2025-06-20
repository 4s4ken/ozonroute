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
	"gitlab.ozon.dev/gleb_evseev/homework/pkg/pvz"
	"gitlab.ozon.dev/gleb_evseev/homework/pkg/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"gitlab.ozon.dev/gleb_evseev/homework/internal/mw"
)

const (
	grpcAddress    = "localhost:50051"
	swaggerAddress = "localhost:7002"
)

func main() {
	ctx := context.Background()

	database := storage.NewPgxPool(ctx, config.DBConnString)
	defer database.Close()

	dbStorage := db.New(database)
	handler := app.NewService(dbStorage, func() time.Time { return time.Now().UTC() })

	go func() {
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

		log.Println("gRPC server started on", grpcAddress)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve gRPC: %v", err)
		}
	}()

	go func() {
		mux := chi.NewMux()

		gw := runtime.NewServeMux()
		err := usecase.RegisterOrdersServiceHandlerFromEndpoint(
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

		mux.Handle("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL("/swagger.json"),
		))

		log.Println("HTTP Gateway + Swagger started on", swaggerAddress)
		if err := http.ListenAndServe(swaggerAddress, mux); err != nil {
			log.Fatalf("failed to serve HTTP: %v", err)
		}
	}()

	select {}
}
