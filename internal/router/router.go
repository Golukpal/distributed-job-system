package router

import (
    "github.com/gofiber/fiber/v2"
    "github.com/jackc/pgx/v5/pgxpool"

    "github.com/Golukpal/job-distributor/internal/handler"
    "github.com/Golukpal/job-distributor/internal/repository"
)

func SetupRoutes(app *fiber.App, db *pgxpool.Pool) {
    healthHandler := handler.NewHealthHandler(db)

    jobRepo := repository.NewJobRepository(db)
    jobHandler := handler.NewJobHandler(jobRepo)

    app.Get("/health", healthHandler.Check)

    api := app.Group("/api/v1")

    api.Get("/health", healthHandler.Check)

    jobs := api.Group("/jobs")

    jobs.Post("/", jobHandler.CreateJob)
    jobs.Get("/:id", jobHandler.GetJob)
}