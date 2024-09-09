package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/gabrielmoura/raspController/configs"
	"github.com/gabrielmoura/raspController/infra/db"
	"github.com/gabrielmoura/raspController/infra/gpio"
	"github.com/gabrielmoura/raspController/infra/routes"
	"github.com/gabrielmoura/raspController/internal/install"
	"github.com/gabrielmoura/raspController/pkg/mdns"
	"github.com/gofiber/fiber/v2"
)

func main() {
	// Parsing the install flag
	installFlag := flag.Bool("install", false, "Run the installation process")
	flag.Parse()

	if *installFlag {
		install.Install()
		return
	}

	// Start the main application
	if err := run(); err != nil {
		log.Fatalf("Application failed: %v", err)
	}
}

// run initializes and starts the application
func run() error {
	// Create a new context
	ctx := context.Background()

	// Load configuration
	if err := configs.LoadConfig(); err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		Prefork:       false,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  configs.Conf.AppName,
		AppName:       configs.Conf.AppName,
	})

	// Initialize routes
	routes.InitializeRoutes(app)

	// Initialize the database
	if err := db.Initialize(ctx); err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	// Initialize GPIO
	if err := gpio.Initialize(ctx); err != nil {
		fmt.Println("failed to initialize GPIO: %w", err)
	}

	// Set mDNS
	if err := mdns.SetDNS(configs.Conf.AppName, configs.Conf.Port); err != nil {
		log.Println("Warning: Failed to set mDNS:", err)
	}

	// Start Fiber server
	address := fmt.Sprintf(":%d", configs.Conf.Port)
	log.Printf("Starting server on %s", address)
	return app.Listen(address)
}
