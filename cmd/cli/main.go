package main

import (
	"context"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/infra/cli"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/infra/storage"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/usecase"
	"log"
	"time"
)

func main() {
	ctx := context.Background()
	js, err := storage.NewJsonStorage("db.json", func() time.Time { return time.Now().UTC() })
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}

	handler := usecase.New(js)

	fileHandler := cli.NewFileHandler(handler)

	err = fileHandler.Run(ctx)
	if err != nil {
		log.Fatalf("")
	}
}
