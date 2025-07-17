package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/swaggo/http-swagger"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/app"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/config"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
	hl "gitlab.ozon.dev/gleb_evseev/homework/internal/infra/grpc"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/kafka/producer"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/kafka/worker"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/metrics"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/storage/db"
	ob "gitlab.ozon.dev/gleb_evseev/homework/internal/storage/outbox"
	"gitlab.ozon.dev/gleb_evseev/homework/pkg/cache"
	usecase "gitlab.ozon.dev/gleb_evseev/homework/pkg/pvz"
	"gitlab.ozon.dev/gleb_evseev/homework/pkg/storage"
	"gitlab.ozon.dev/gleb_evseev/homework/pkg/tracer"
	"gitlab.ozon.dev/gleb_evseev/homework/pkg/workerpool"
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
	kafkaBroker    = "localhost:9092"
	kafkaTopic     = "pvz.events-log"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	shutdown := tracer.InitTracer(ctx, "pvz-app")
	defer shutdown()

	database := storage.NewPgxPool(ctx, config.DBConnString)
	dbStorage := db.New(database)
	defer func() {
		if err := dbStorage.Close(); err != nil {
			log.Printf("failed to close storage: %v", err)
		}
	}()

	wp := workerpool.NewWorkerPool(ctx, 10)
	orderCache := cache.New[string, *pvz.Order](config.CacheCfg)
	outboxRepo := ob.NewOutboxRepo(database)
	kafkaProducer, err := producer.NewProducer([]string{kafkaBroker}, kafkaTopic)
	if err != nil {
		log.Fatalf("failed to create kafka producer: %v", err)
	}
	outboxWorker := worker.NewOutboxWorker(outboxRepo, kafkaProducer, 10)
	go outboxWorker.Run(ctx)

	handler := app.NewService(dbStorage, outboxRepo, func() time.Time { return time.Now().UTC() }, wp, orderCache)

	lis, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			mw.LoggingInterceptor,
			mw.ValidateInterceptor,
			mw.TraceInterceptor,
			mw.RateLimiterInterceptor(
				limiter.New(
					memory.NewStore(),
					limiter.Rate{Period: 200 * time.Millisecond, Limit: 2},
				),
			),
			mw.ErrorMappingInterceptor(),
			mw.WorkerPoolInterceptor(wp),
		),
	)

	usecase.RegisterOrdersServiceServer(grpcServer, hl.NewOrdersServer(handler))
	go func() {
		log.Println("gRPC server started on", grpcAddress)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve gRPC: %v", err)
		}
	}()

	metricsMW := mw.NewMetricsMW()
	metricsMW.Register(prometheus.DefaultRegisterer)
	metrics.Register()

	mux := chi.NewMux()

	mux.Use(metricsMW.Handler)
	mux.Handle("/metrics", promhttp.Handler())

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

	cancel()
	wp.Shutdown()

	grpcServer.GracefulStop()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer shutdownCancel()
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("Ошибка при остановке HTTP сервера: %v", err)
	}

	if err := dbStorage.Close(); err != nil {
		log.Printf("Ошибка при закрытии БД: %v", err)
	}

}
