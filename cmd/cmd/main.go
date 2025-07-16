package main

import (
	"context"
	"github.com/patrickmn/go-cache"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/app"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/config"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/infra/cli"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/storage/json"
	"gitlab.ozon.dev/gleb_evseev/homework/pkg/workerpool"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	js, err := json.NewJsonStorage(config.DBPath, func() time.Time { return time.Now().UTC() })
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}

	wp := workerpool.NewWorkerPool(ctx, 5)
	defer wp.Shutdown()
	orderCache := *cache.New(1*time.Minute, 10*time.Second)
	outboxRepo := json.New(config.OutboxPath)
	handler := app.NewService(js, outboxRepo, func() time.Time { return time.Now().UTC() }, wp, orderCache)
	fileHandler := cli.NewFileHandler(handler, wp)

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		log.Println("Graceful shutdown")
		cancel()
		wp.Shutdown()
	}()

	if err := fileHandler.Run(ctx, wp); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
