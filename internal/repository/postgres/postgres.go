package postgres

import (
	"github.com/DimaGitHahahab/ozon-fintech-posts/internal/repository"
	"github.com/DimaGitHahahab/ozon-fintech-posts/internal/repository/postgres/queries"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresRepo struct {
	*queries.Queries
	pool *pgxpool.Pool
}

func New(pgxPool *pgxpool.Pool) repository.Repository {
	r := &postgresRepo{
		Queries: queries.New(pgxPool),
		pool:    pgxPool,
	}

	return r
}
