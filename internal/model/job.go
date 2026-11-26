package model

import (
    "encoding/json"
    "time"
)

type JobStatus string

const (
    JobPending    JobStatus = "pending"
    JobProcessing JobStatus = "processing"
    JobCompleted  JobStatus = "completed"
    JobFailed     JobStatus = "failed"
    JobCancelled  JobStatus = "cancelled"
)

type Job struct {
    ID             string          `json:"id"`
    Type           string          `json:"type"`
    Payload        json.RawMessage `json:"payload"`
    Status         JobStatus       `json:"status"`
    Attempts       int             `json:"attempts"`
    MaxAttempts    int             `json:"max_attempts"`
    ScheduledAt    time.Time       `json:"scheduled_at"`
    StartedAt      *time.Time      `json:"started_at,omitempty"`
    CompletedAt    *time.Time      `json:"completed_at,omitempty"`
    FailedAt       *time.Time      `json:"failed_at,omitempty"`
    LastError      *string         `json:"last_error,omitempty"`
    IdempotencyKey *string         `json:"-"`
    CreatedAt      time.Time       `json:"created_at"`
    UpdatedAt      time.Time       `json:"updated_at"`
}