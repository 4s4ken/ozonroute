package main

import (
	"context"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/app"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/config"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/infra/cli"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/storage/json"
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

	wp := app.NewWorkerPool(ctx, 5)
	defer wp.Shutdown()

	handler := app.NewService(js, func() time.Time { return time.Now().UTC() }, wp)
	fileHandler := cli.NewFileHandler(handler, wp)

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		log.Println("Graceful shutdown")
		cancel()
		wp.Shutdown()
	}()

	if err := fileHandler.Run(ctx); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
