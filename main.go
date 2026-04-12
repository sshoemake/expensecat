package main

import (
	"expensecat/internal/cli"
	"expensecat/internal/storage"
	"log"
)

func main() {
	store, err := storage.InitializeStorage()
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}
	defer store.Close()

	if err := cli.RunApp(store); err != nil {
		log.Fatalf("Failed to run application: %v", err)
	}
}
