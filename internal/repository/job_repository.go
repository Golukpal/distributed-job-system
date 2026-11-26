package repository

import (
    "context"
    "encoding/json"
    "errors"
    "fmt"

    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"

    "github.com/Golukpal/job-distributor/internal/model"
)

var ErrJobNotFound = errors.New("job not found")

type JobRepository struct {
    DB *pgxpool.Pool
}

func NewJobRepository(db *pgxpool.Pool) *JobRepository {
    return &JobRepository{
        DB: db,
    }
}

func (r *JobRepository) CreateJob(
    ctx context.Context,
    job *model.Job,
) (*model.Job, error) {

    query := `
        INSERT INTO jobs (
            type,
            payload,
            status,
            max_attempts,
            idempotency_key
        )
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT (idempotency_key)
        WHERE idempotency_key IS NOT NULL
        DO NOTHING
        RETURNING
            id,
            type,
            payload,
            status,
            attempts,
            max_attempts,
            scheduled_at,
            started_at,
            completed_at,
            failed_at,
            last_error,
            idempotency_key,
            created_at,
            updated_at
    `

    payload, err := json.Marshal(job.Payload)
    if err != nil {
        return nil, fmt.Errorf("failed to encode payload: %w", err)
    }

    result, err := scanJob(
        r.DB.QueryRow(
            ctx,
            query,
            job.Type,
            payload,
            model.JobPending,
            job.MaxAttempts,
            job.IdempotencyKey,
        ),
    )

    if err == nil {
        return result, nil
    }

    if !errors.Is(err, pgx.ErrNoRows) {
        return nil, fmt.Errorf("failed to create job: %w", err)
    }

    // An existing idempotency key means the original
    // job should be returned instead of creating a duplicate.
    if job.IdempotencyKey != nil {
        existing, getErr := r.GetJobByIdempotencyKey(
            ctx,
            *job.IdempotencyKey,
        )

        if getErr != nil {
            return nil, getErr
        }

        return existing, nil
    }

    return nil, fmt.Errorf("failed to create job: %w", err)
}

func (r *JobRepository) GetJobByID(
    ctx context.Context,
    id string,
) (*model.Job, error) {

    query := `
        SELECT
            id,
            type,
            payload,
            status,
            attempts,
            max_attempts,
            scheduled_at,
            started_at,
            completed_at,
            failed_at,
            last_error,
            idempotency_key,
            created_at,
            updated_at
        FROM jobs
        WHERE id::text = $1
    `

    job, err := scanJob(r.DB.QueryRow(ctx, query, id))

    if errors.Is(err, pgx.ErrNoRows) {
        return nil, ErrJobNotFound
    }

    if err != nil {
        return nil, fmt.Errorf("failed to retrieve job: %w", err)
    }

    return job, nil
}

func (r *JobRepository) GetJobByIdempotencyKey(
    ctx context.Context,
    key string,
) (*model.Job, error) {

    query := `
        SELECT
            id,
            type,
            payload,
            status,
            attempts,
            max_attempts,
            scheduled_at,
            started_at,
            completed_at,
            failed_at,
            last_error,
            idempotency_key,
            created_at,
            updated_at
        FROM jobs
        WHERE idempotency_key = $1
    `

    job, err := scanJob(r.DB.QueryRow(ctx, query, key))

    if errors.Is(err, pgx.ErrNoRows) {
        return nil, ErrJobNotFound
    }

    if err != nil {
        return nil, fmt.Errorf("failed to retrieve job: %w", err)
    }

    return job, nil
}

func scanJob(row pgx.Row) (*model.Job, error) {
    var job model.Job

    err := row.Scan(
        &job.ID,
        &job.Type,
        &job.Payload,
        &job.Status,
        &job.Attempts,
        &job.MaxAttempts,
        &job.ScheduledAt,
        &job.StartedAt,
        &job.CompletedAt,
        &job.FailedAt,
        &job.LastError,
        &job.IdempotencyKey,
        &job.CreatedAt,
        &job.UpdatedAt,
    )

    if err != nil {
        return nil, err
    }

    return &job, nil
}