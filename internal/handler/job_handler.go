package handler

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/Golukpal/job-distributor/internal/model"
	"github.com/Golukpal/job-distributor/internal/repository"
)

type JobHandler struct {
	JobRepo *repository.JobRepository
}

func NewJobHandler(repo *repository.JobRepository) *JobHandler {
	return &JobHandler{
		JobRepo: repo,
	}
}

type CreateJobRequest struct {
	Type        string          `json:"type"`
	Payload     json.RawMessage `json:"payload"`
	MaxAttempts int             `json:"max_attempts"`
}

func (h *JobHandler) CreateJob(c *fiber.Ctx) error {
	var req CreateJobRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	req.Type = strings.TrimSpace(req.Type)

	if req.Type == "" || len(req.Type) > 100 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "job type must contain between 1 and 100 characters",
		})
	}

	if len(req.Payload) == 0 || !json.Valid(req.Payload) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "valid JSON payload is required",
		})
	}

	if len(req.Payload) > 1024*1024 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "payload cannot exceed 1 MB",
		})
	}

	if req.MaxAttempts == 0 {
		req.MaxAttempts = 3
	}

	if req.MaxAttempts < 1 || req.MaxAttempts > 10 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "max_attempts must be between 1 and 10",
		})
	}

	idempotencyKey := strings.TrimSpace(
		c.Get("Idempotency-Key"),
	)

	if len(idempotencyKey) > 255 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "idempotency key cannot exceed 255 characters",
		})
	}

	var key *string

	if idempotencyKey != "" {
		key = &idempotencyKey
	}

	job := &model.Job{
		Type:           req.Type,
		Payload:        req.Payload,
		Status:         model.JobPending,
		MaxAttempts:    req.MaxAttempts,
		IdempotencyKey: key,
	}

	ctx, cancel := timeBoundContext(c)
	defer cancel()

	createdJob, err := h.JobRepo.CreateJob(ctx, job)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create job",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(createdJob)
}

func (h *JobHandler) GetJob(c *fiber.Ctx) error {
	id := c.Params("id")

	if strings.TrimSpace(id) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "job ID is required",
		})
	}

	ctx, cancel := timeBoundContext(c)
	defer cancel()

	job, err := h.JobRepo.GetJobByID(ctx, id)

	if errors.Is(err, repository.ErrJobNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "job not found",
		})
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to retrieve job",
		})
	}

	return c.Status(fiber.StatusOK).JSON(job)
}

func timeBoundContext(c *fiber.Ctx) (context.Context, context.CancelFunc) {
	return context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
}
