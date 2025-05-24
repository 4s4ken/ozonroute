package main

import (
	"bufio"
	"fmt"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/infra"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/storage/json_store"
	"log"
	"os"
	"strings"
)

func main() {
	js, err := json_store.NewJsonStorage("db.json")
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}

	st := infra.StorageService{
		Store: js,
	}
	rootCmd := infra.InitCommands(&st)

	fmt.Println("Welcome to the order pick-up point(PVZ) simulation. To view all available commands, type 'help'.")

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print(">")
		if !scanner.Scan() {
			break
		}
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		rootCmd.SetArgs(parts)

		if err := rootCmd.Execute(); err != nil {
			mapped := json_store.MapError(err)
			_, _ = fmt.Fprintf(os.Stderr, "Command failed: %s\n", mapped)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Input error: %s", err)
	}
}
