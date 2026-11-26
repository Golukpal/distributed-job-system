CREATE TABLE IF NOT EXISTS jobs

(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    type VARCHAR(100) NOT NULL,

    payload JSONB NOT NULL,

    status VARCHAR(20) NOT NULL DEFAULT 'pending',

    attempts INTEGER NOT NULL DEFAULT 0,

    max_attempts INTEGER NOT NULL DEFAULT 3,

    scheduled_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    started_at TIMESTAMPTZ ,

    completed_at TIMESTAMPTZ,

    failed_at TIMESTAMPTZ,

    last_error TEXT,

    CONSTRAINT valid_job_status
        CHECK(
            status IN (
                'pending',
                'processing',
                'completed',
                'failed'
            )
        ),
    CONSTRAINT valid_attempts
        CHECK (
            attempts >= 0 AND max_attempts BETWEEN 1 AND 10
        )

);

    CREATE UNIQUE INDEX IF NOT EXISTS idx_jobs_idempotency_key
    ON jobs (idempotency_key)
    WHERE idempotency_key IS NOT NULL;

    CREATE INDEX IF NOT EXISTS idx_jobs_status_scheduled
    ON jobs (status, scheduled_at);

    CREATE INDEX IF NOT EXISTS idx_jobs_created_at
    ON jobs (created_at DESC);
