package queries

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

// Queries implements the SQL queries for postgresRepo
type Queries struct {
	pool *pgxpool.Pool
}

// New creates a new instance of Queries to be embedded in postgresRepo
func New(pgxPool *pgxpool.Pool) *Queries {
	return &Queries{pool: pgxPool}
}
