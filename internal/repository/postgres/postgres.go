package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/DimaGitHahahab/ozon-fintech-posts/internal/repository"
	"github.com/DimaGitHahahab/ozon-fintech-posts/internal/repository/postgres/queries"
	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5/pgxpool"
)

// postgresRepo is a PostgreSQL implementation of repository.Repository
type postgresRepo struct {
	*queries.Queries // SQL queries
	pool             *pgxpool.Pool
}

// New creates a new instance of the repository based on connection pool
func New(pgxPool *pgxpool.Pool) repository.Repository {
	r := &postgresRepo{
		Queries: queries.New(pgxPool),
		pool:    pgxPool,
	}

	return r
}

// SetupPgxPool creates a new pool of connections
func SetupPgxPool(ctx context.Context, DbURL string) (*pgxpool.Pool, error) {
	pgxConfig, err := pgxpool.ParseConfig(DbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse pgx config: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create new pgx pool with config: %w", err)
	}

	return pool, nil
}

// ProcessMigration runs the migration up
func ProcessMigration(migrationURL string, dbSource string) error {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		return fmt.Errorf("failed to create new migration: %w", err)
	}

	if err = migration.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to migrate up: %w", err)
	}
	defer migration.Close()

	return nil
}
