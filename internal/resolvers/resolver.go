package resolvers

import (
	"context"

	"github.com/DimaGitHahahab/ozon-fintech-posts/internal/repository"
	"github.com/DimaGitHahahab/ozon-fintech-posts/internal/schema"
)

type Resolver struct {
	repo repository.Repository
}

func (r Resolver) GetPosts() (any, error) {
	// TODO
	panic("implement me")
}

func (r Resolver) GetPost(ctx context.Context, args schema.PostArgs) (any, error) {
	// TODO
	panic("implement me")
}

func (r Resolver) GetCommentsByPost(ctx context.Context, args schema.GetCommentsArgs) (any, error) {
	// TODO
	panic("implement me")
}

func (r Resolver) GetCommentsByParent(ctx context.Context, args schema.GetCommentsArgs) (any, error) {
	// TODO
	panic("implement me")
}

func (r Resolver) CreatePost(ctx context.Context, args schema.CreatePostArgs) (any, error) {
	// TODO
	panic("implement me")
}

func (r Resolver) CreateComment(ctx context.Context, args schema.CreateCommentArgs) (any, error) {
	// TODO
	panic("implement me")
}

func (r Resolver) DisableComments(ctx context.Context, args schema.DisableCommentsArgs) (any, error) {
	// TODO
	panic("implement me")
}
