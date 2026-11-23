package main

import (
	"log"
	"github.com/Golukpal/job-distributor/internal/config"
	"github.com/Golukpal/job-distributor/internal/database"
	"github.com/Golukpal/job-distributor/internal/router"
	"github.com/gofiber/fiber/v2"
)

func main() {
	cfg := config.Load()

	db, err := database.ConnectDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer db.Close()

	app := fiber.New()

	router.SetupRoutes(app, db)

	log.Printf("%s starting on port %s", cfg.AppName, cfg.AppPort)

	if err := app.Listen(":" + cfg.AppPort); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
