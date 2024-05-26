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
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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
		log.Println("Processing migration...")
		if err := postgres.ProcessMigration(cfg.MigrationPath, cfg.DbURL); err != nil {
			return nil, fmt.Errorf("failed to process migration: %w", err)
		}
		log.Println("Migration is successful")

		log.Println("Setting up pgx pool...")
		pool, err := postgres.SetupPgxPool(ctx, cfg.DbURL)
		if err != nil {
			return nil, fmt.Errorf("failed to setup pgx pool: %w", err)
		}
		log.Println("Pgx pool is set up successfully")

		return postgres.New(pool), nil
	default:
		return nil, fmt.Errorf("unknown repository type: %s", cfg.Repository)
	}
}
