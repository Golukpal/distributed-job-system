package database

import (
    "context"
    "embed"
    "fmt"

    "github.com/jackc/pgx/v5/pgxpool"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

func RunMigrations(ctx context.Context, db *pgxpool.Pool) error {

    tx, err := db.Begin(ctx)
    if err != nil {
        return fmt.Errorf("failed to start migration transaction: %w", err)
    }

    defer tx.Rollback(ctx)

    // Prevent multiple application instances from
    // running the same migration simultaneously.
    _, err = tx.Exec(ctx,
        `SELECT pg_advisory_xact_lock(741852963)`,
    )
    if err != nil {
        return fmt.Errorf("failed to acquire migration lock: %w", err)
    }

    _, err = tx.Exec(ctx, `
        CREATE TABLE IF NOT EXISTS schema_migrations (
            version VARCHAR(255) PRIMARY KEY,
            applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
        )
    `)
    if err != nil {
        return fmt.Errorf("failed to create migrations table: %w", err)
    }

    const version = "001_create_jobs"

    var exists bool

    err = tx.QueryRow(ctx, `
        SELECT EXISTS (
            SELECT 1
            FROM schema_migrations
            WHERE version = $1
        )
    `, version).Scan(&exists)

    if err != nil {
        return fmt.Errorf("failed to check migration version: %w", err)
    }

    if exists {
        return tx.Commit(ctx)
    }

    sqlFile, err := migrationFiles.ReadFile(
        "migrations/001_create_jobs.sql",
    )
    if err != nil {
        return fmt.Errorf("failed to read migration file: %w", err)
    }

    _, err = tx.Exec(ctx, string(sqlFile))
    if err != nil {
        return fmt.Errorf("failed to execute migration: %w", err)
    }

    _, err = tx.Exec(ctx, `
        INSERT INTO schema_migrations (version)
        VALUES ($1)
    `, version)

    if err != nil {
        return fmt.Errorf("failed to record migration: %w", err)
    }

    if err := tx.Commit(ctx); err != nil {
        return fmt.Errorf("failed to commit migration: %w", err)
    }

    return nil
}