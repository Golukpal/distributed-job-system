package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type HealthHandler struct {
	DB *pgxpool.Pool
}

func NewHealthHandler(db *pgxpool.Pool) *HealthHandler {
	return &HealthHandler{
		DB: db,
	}
}

func (h *HealthHandler) Check(c *fiber.Ctx) error {
	if err := h.DB.Ping(c.Context()); err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(
			fiber.Map{
				"status":   "unhealthy",
				"database": "disconnected",
			},
		)
	}

	return c.Status(fiber.StatusOK).JSON(
		fiber.Map{
			 "status": "healthy",
            "database": "connected",
		},
	)
}
