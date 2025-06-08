package main

import (
	"context"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/app"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/config"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/infra/cli"
	"log"
	"time"
)

func main() {
	ctx := context.Background()
	js, err := app.NewJsonStorage(config.DBPath, func() time.Time { return time.Now().UTC() })
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}

	handler := app.New(js)

	fileHandler := cli.NewFileHandler(handler)

	err = fileHandler.Run(ctx)
	if err != nil {
		log.Fatalf("Error: Failed to run command: %v", err)
	}
}
