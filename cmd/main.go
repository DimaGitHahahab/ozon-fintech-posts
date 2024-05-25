package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/DimaGitHahahab/ozon-fintech-posts/internal/domain"
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

	post, err := repo.CreatePost(ctx, &domain.Post{
		Title:            "My first post!",
		Content:          "lorem ipsum",
		AuthorID:         15,
		CreatedAt:        time.Now(),
		CommentsDisabled: false,
	})
	fmt.Println(post)
	fmt.Println()

	posts, err := repo.GetPosts(ctx)
	if err != nil {
		log.Fatal(err)
	}

	for _, p := range posts {
		fmt.Println(p)
	}
	fmt.Println()

	com, err := repo.CreateComment(ctx, &domain.Comment{
		PostID:    1,
		ParentID:  nil,
		AuthorID:  10,
		Content:   "nice post, i like it!",
		CreatedAt: time.Now(),
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(com)
	fmt.Println()

	comments, err := repo.GetCommentsByPost(ctx, 1, 10, 0)
	if err != nil {
		log.Fatal(err)
	}

	for _, c := range comments {
		fmt.Println(c)
	}
}
