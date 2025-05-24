package main

import (
	"bufio"
	"fmt"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/infra/cli"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/infra/storage"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/usecase"
	"log"
	"os"
	"strings"
)

func main() {
	js, err := storage.NewJsonStorage("db.json")
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}

	handler := usecase.New(js)

	fileHandler := cli.NewFileHandler(handler)
	err := fileHandler.Run()
	if err != nil {
		log.Fatalf("")
	}
}
