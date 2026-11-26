package main

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/Golukpal/job-distributor/internal/config"
	"github.com/Golukpal/job-distributor/internal/database"
	"github.com/Golukpal/job-distributor/internal/router"
)

func main() {
	cfg := config.Load()

	db, err := database.ConnectDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer db.Close()

	migrationCtx, cancel := context.WithTimeout(
		context.Background(),
		30*time.Second,
	)

	err = database.RunMigrations(migrationCtx, db)
	cancel()

	if err != nil {
		log.Fatalf("Database migration failed: %v", err)
	}

	log.Println("Database migrations completed")

	app := fiber.New(fiber.Config{
		AppName:   cfg.AppName,
		BodyLimit: 2 * 1024 * 1024,
	})

	router.SetupRoutes(app, db)

	log.Printf("%s starting on port %s", cfg.AppName, cfg.AppPort)

	if err := app.Listen(":" + cfg.AppPort); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
