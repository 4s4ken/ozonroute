package main

import (
	"bufio"
	"fmt"
	"gitlab.ozon.dev/gleb_evseev/homework/homework-1/internal/storage/json_store"
	"gitlab.ozon.dev/gleb_evseev/homework/homework-1/internal/usecase"
	"log"
	"os"
	"strings"
)

func main() {
	js, err := json_store.NewJsonStorage("db.json")
	if err != nil {
		log.Fatalf("cannot init json storage: %s", err)
	}

	st := usecase.StorageService{
		Store: js,
	}
	rootCmd := usecase.InitCommands(&st)

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
			fmt.Fprintf(os.Stderr, "Command failed, Input Error: %s\n", err)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error: %s", err)
	}
}
