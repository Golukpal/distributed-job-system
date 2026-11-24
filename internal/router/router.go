package router

import (
	"github.com/Golukpal/job-distributor/internal/handler"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupRoutes(app *fiber.App, db *pgxpool.Pool) {
	healthHandler := handler.NewHealthHandler(db)
	// app.Get("/health", healthHandler.Check)
	api := app.Group("/api/v1")

	api.Get("/health", healthHandler.Check) 
}
