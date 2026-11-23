package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectDB(databseURL string) (*pgxpool.Pool, error) {
	if databseURL == "" {
		return nil, errors.New("database URL not found")
	}

	config, err := pgxpool.ParseConfig(databseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse db config: %w", err)
	}

	config.MaxConns = 20
	config.MinConns = 2
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil{
		return nil, fmt.Errorf("failed to create db pool: %w", err)
	}

	if err := pool.Ping(ctx); err!= nil {
		pool.Close()
		return nil, fmt.Errorf(" failed to connect to db %w", err)
	}

	return pool, nil



}
