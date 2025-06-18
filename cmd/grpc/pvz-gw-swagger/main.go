package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/swaggo/http-swagger"
	desc "gitlab.ozon.dev/gleb_evseev/homework/pkg/pvz"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	grpcAddress    = "localhost:50051"
	swaggerAddress = "localhost:7002"
)

func main() {
	ctx := context.Background()
	mux := chi.NewMux()

	gatewayMux := runtime.NewServeMux()
	err := desc.RegisterOrdersServiceHandlerFromEndpoint(
		ctx, gatewayMux, grpcAddress,
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	)
	if err != nil {
		log.Fatalf("failed to register gRPC gateway: %v", err)
	}

	mux.Mount("/", gatewayMux)

	mux.Get("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		b, err := os.ReadFile("./pkg/pvz/pvz.swagger.json")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write(b); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	})

	mux.Handle("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger.json"),
	))

	log.Printf("HTTP server with Swagger running on %s", swaggerAddress)
	if err := http.ListenAndServe(swaggerAddress, mux); err != nil {
		log.Fatalf("failed to serve HTTP: %v", err)
	}
}
