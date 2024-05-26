package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/DimaGitHahahab/ozon-fintech-posts/internal/repository"
	"github.com/DimaGitHahahab/ozon-fintech-posts/internal/repository/in_memory"
	"github.com/DimaGitHahahab/ozon-fintech-posts/internal/repository/postgres"
	"github.com/DimaGitHahahab/ozon-fintech-posts/internal/resolvers"
	"github.com/DimaGitHahahab/ozon-fintech-posts/internal/schema"
	"github.com/DimaGitHahahab/ozon-fintech-posts/internal/server"
	"github.com/DimaGitHahahab/ozon-fintech-posts/pkg/config"
	"github.com/DimaGitHahahab/ozon-fintech-posts/pkg/signal"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	config  *config.Config
	sigQuit chan os.Signal
	srv     *server.Server
}

func New(cfg *config.Config) *App {
	ctx := context.Background()

	repo, err := createRepository(ctx, cfg)
	if err != nil {
		log.Fatal("Failed to create repository: ", err)
	}

	resolver := resolvers.NewResolver(repo)

	sch, err := schema.NewSchema(resolver)
	if err != nil {
		log.Fatal("Failed to create new GraphQL schema: ", err)
	}

	srv := server.NewServer(&sch)

	return &App{
		config:  cfg,
		srv:     srv,
		sigQuit: signal.GetShutdownChannel(),
	}
}

func (a *App) Run() {
	go func() {
		log.Println("Starting server on port", a.config.HTTPPort)
		if err := a.srv.Run(a.config.HTTPPort); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalln("Failed to start server: ", err)
		}
	}()

	<-a.sigQuit
	log.Println("Gracefully shutting down server")

	if err := a.srv.Shutdown(context.Background()); err != nil {
		log.Fatalln("Failed to shutdown the server gracefully: ", err)
	}

	log.Println("Server shutdown is successful")
}

func createRepository(ctx context.Context, cfg *config.Config) (repository.Repository, error) {
	switch cfg.Repository {
	case "IN_MEMORY":
		return in_memory.New(), nil
	case "POSTGRES":
		if err := processMigration(cfg.MigrationPath, cfg.DbURL); err != nil {
			return nil, fmt.Errorf("failed to process migration: %w", err)
		}

		pool, err := setupPgxPool(ctx, cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to setup pgx pool: %w", err)
		}
		return postgres.New(pool), nil
	default:
		return nil, fmt.Errorf("unknown repository type: %s", cfg.Repository)
	}
}

func setupPgxPool(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	log.Println("Setting up pgx pool...")

	pgxConfig, err := pgxpool.ParseConfig(cfg.DbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse pgx config: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create new pgx pool with config: %w", err)
	}

	log.Println("Pgx pool initialized successfully")
	return pool, nil
}

func processMigration(migrationURL string, dbSource string) error {
	log.Println("Processing migration...")

	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		return fmt.Errorf("failed to create new migration: %w", err)
	}

	if err = migration.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to migrate up: %w", err)
	}
	defer migration.Close()

	log.Println("Migration successful")
	return nil
}
