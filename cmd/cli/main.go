package main

import (
	"context"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/infra/cli"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/infra/storage"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/usecase"
	"log"
)

func main() {
	ctx := context.Background()
	js, err := storage.NewJsonStorage("some path")
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
