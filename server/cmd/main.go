package main

import (
	"flag"
	"log"
	"os"

	"sol_privacy/internal/cli"
	"sol_privacy/internal/server"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Define flags
	serverMode := flag.Bool("server", false, "Run as HTTP API server instead of CLI")
	port := flag.String("port", "8080", "Port to run server on (only used with --server)")
	flag.Parse()

	if *serverMode {
		runServer(*port)
	} else {
		runCLI()
	}
}

func runCLI() {
	if err := cli.Run(); err != nil {
		log.Fatal(err)
	}
}

func runServer(port string) {
	apiKey := os.Getenv("SHADOWPAY_API_KEY")
	if apiKey == "" {
		log.Fatal("SHADOWPAY_API_KEY environment variable is required")
	}

	if err := server.Run(server.Config{
		APIKey: apiKey,
		Port:   port,
	}); err != nil {
		log.Fatal(err)
	}
}
