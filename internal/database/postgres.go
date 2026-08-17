package database
import (
  	"context"
  	"errors"
  	"fmt"

  	"github.com/jackc/pgx/v5/pgxpool"
  )

  func OpenPostgres(
  	ctx context.Context,
  	databaseURL string,
  ) (*pgxpool.Pool, error) {
  	if databaseURL == "" {
  		return nil, errors.New("DATABASE_URL is required")
  	}

  	poolConfig, err := pgxpool.ParseConfig(databaseURL)
  	if err != nil {
  		return nil, fmt.Errorf("parse PostgreSQL configuration: %w", err)
  	}

  	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
  	if err != nil {
  		return nil, fmt.Errorf("create PostgreSQL pool: %w", err)
  	}

  	if err := pool.Ping(ctx); err != nil {
  		pool.Close()
  		return nil, fmt.Errorf("ping PostgreSQL: %w", err)
  	}

  	return pool, nil
  }
