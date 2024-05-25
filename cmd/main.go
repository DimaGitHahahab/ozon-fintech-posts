package main

import (
	"context"
	"fmt"
	"log"

	"github.com/DimaGitHahahab/ozon-fintech-posts/internal/repository/postgres"
	"github.com/DimaGitHahahab/ozon-fintech-posts/pkg/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config: ", err)
	}

	pgxConfig, err := pgxpool.ParseConfig(cfg.DbURL)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	pool, err := pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		log.Fatal(err)
	}

	repo := postgres.New(pool)

	post, err := repo.GetPost(ctx, 1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(post)
	comments, err := repo.GetCommentsByPost(ctx, 1, 5, 0)
	if err != nil {
		log.Fatal(err)
	}
	for _, c := range comments {
		fmt.Println(c)
	}
	fmt.Println()

	post, err = repo.GetPost(ctx, 2)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(post)
	comments, err = repo.GetCommentsByPost(ctx, 2, 5, 0)
	if err != nil {
		log.Fatal(err)
	}
	for _, c := range comments {
		fmt.Println(c)
	}
}
