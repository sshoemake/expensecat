package main

import (
	"expensecat/internal/cli"
	"expensecat/internal/storage"
	"expensecat/internal/version"
	"flag"
	"fmt"
	"log"
	"os"
)

// Set via ldflags at build time
var (
	appVersion = "dev"
	commit     = "unknown"
	buildTime  = "unknown"
)

func init() {
	if appVersion != "dev" {
		version.Version = appVersion
	}
	if commit != "unknown" {
		version.Commit = commit
	}
	if buildTime != "unknown" {
		version.BuildTime = buildTime
	}
	if version.Version != "dev" && appVersion == "dev" {
		appVersion = version.Version
	}
	if version.Commit != "unknown" && commit == "unknown" {
		commit = version.Commit
	}
	if version.BuildTime != "unknown" && buildTime == "unknown" {
		buildTime = version.BuildTime
	}
}

func main() {
	showVersion := flag.Bool("version", false, "Show version information")
	flag.Parse()

	if *showVersion {
		fmt.Printf("expensecat %s\n", appVersion)
		fmt.Printf("Commit: %s\n", commit)
		fmt.Printf("Built:  %s\n", buildTime)
		os.Exit(0)
	}

	store, err := storage.InitializeStorage()
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}
	defer store.Close()

	if err := cli.RunApp(store); err != nil {
		log.Fatalf("Failed to run application: %v", err)
	}
}
