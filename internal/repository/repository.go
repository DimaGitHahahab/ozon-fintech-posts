package repository

import (
	"context"

	"github.com/DimaGitHahahab/ozon-fintech-posts/internal/domain"
)

type Repository interface {
	GetPosts(ctx context.Context) ([]*domain.Post, error)
	GetPost(ctx context.Context, id int) (*domain.Post, error)
	CreatePost(ctx context.Context, post *domain.Post) (*domain.Post, error)
	CreateComment(ctx context.Context, comment *domain.Comment) (*domain.Comment, error)
	GetCommentsByPost(ctx context.Context, postID int, limit, offset int) ([]*domain.Comment, error)
	GetCommentsByParent(ctx context.Context, parentId int, limit, offset int) ([]*domain.Comment, error)
	DisableComments(ctx context.Context, postID int) error
}
