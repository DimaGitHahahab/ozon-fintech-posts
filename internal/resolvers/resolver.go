package resolvers

import (
	"context"
	"fmt"
	"time"

	"github.com/DimaGitHahahab/ozon-fintech-posts/internal/domain"
	"github.com/DimaGitHahahab/ozon-fintech-posts/internal/repository"
)

type Resolver struct {
	repo repository.Repository
}

func NewResolver(repo repository.Repository) *Resolver {
	return &Resolver{repo: repo}
}

func (r *Resolver) GetPosts(ctx context.Context) (any, error) {
	posts, err := r.repo.GetPosts(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get posts: %w", err)
	}

	return posts, nil
}

func (r *Resolver) GetPost(ctx context.Context, args PostArgs) (any, error) {
	ok, err := r.repo.ContainsPost(ctx, args.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to check post existence: %w", err)
	}
	if !ok {
		return nil, fmt.Errorf("post with id %d not found", args.ID)
	}

	post, err := r.repo.GetPost(ctx, args.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get post: %w", err)
	}

	return post, nil
}

func (r *Resolver) GetCommentsByPost(ctx context.Context, args GetCommentsArgs) (any, error) {
	ok, err := r.repo.ContainsPost(ctx, args.PostID)
	if err != nil {
		return nil, fmt.Errorf("failed to check post existence: %w", err)
	}

	if !ok {
		return nil, fmt.Errorf("post with id %d not found", args.PostID)
	}

	comments, err := r.repo.GetCommentsByPost(ctx, args.PostID, args.Limit, args.Offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get comments: %w", err)
	}

	return comments, nil
}

func (r *Resolver) GetCommentsByParent(ctx context.Context, args GetCommentsArgs) (any, error) {
	ok, err := r.repo.ContainsComment(ctx, args.ParentID)
	if err != nil {
		return nil, fmt.Errorf("failed to check parent comment existence: %w", err)
	}

	if !ok {
		return nil, fmt.Errorf("comment with id %d not found", args.ParentID)
	}

	comments, err := r.repo.GetCommentsByParent(ctx, args.ParentID, args.Limit, args.Offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get comments: %w", err)
	}

	return comments, nil
}

func (r *Resolver) CreatePost(ctx context.Context, args CreatePostArgs) (any, error) {
	post := &domain.Post{
		Title:     args.Title,
		Content:   args.Content,
		AuthorID:  args.AuthorID,
		CreatedAt: time.Now().UTC(),
	}

	savedPost, err := r.repo.CreatePost(ctx, post)
	if err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	return savedPost, nil
}

func (r *Resolver) CreateComment(ctx context.Context, args CreateCommentArgs) (any, error) {
	if !validateComment(args.Content) {
		return nil, fmt.Errorf("comment is too long")
	}

	comment := &domain.Comment{
		PostID:    args.PostID,
		ParentID:  &args.ParentID,
		AuthorID:  args.AuthorID,
		Content:   args.Content,
		CreatedAt: time.Now(),
	}

	ok, err := r.repo.ContainsPost(ctx, args.PostID)
	if err != nil {
		return nil, fmt.Errorf("failed to check post existence: %w", err)
	}
	if !ok {
		return nil, fmt.Errorf("post with id %d not found", args.PostID)
	}

	post, err := r.repo.GetPost(ctx, args.PostID)
	if err != nil {
		return nil, fmt.Errorf("failed to get post: %w", err)
	}

	if post.CommentsDisabled {
		return nil, fmt.Errorf("comments are disabled for post with id %d", args.PostID)
	}

	savedComment, err := r.repo.CreateComment(ctx, comment)
	if err != nil {
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}

	return savedComment, nil
}

func (r *Resolver) DisableComments(ctx context.Context, args DisableCommentsArgs) (any, error) {
	ok, err := r.repo.ContainsPost(ctx, args.PostID)
	if err != nil {
		return nil, fmt.Errorf("failed to check post existence: %w", err)
	}
	if !ok {
		return nil, fmt.Errorf("post with id %d not found", args.PostID)
	}
	post, err := r.repo.GetPost(ctx, args.PostID)
	if err != nil {
		return nil, fmt.Errorf("failed to get post: %w", err)
	}
	if post.AuthorID != args.AuthorId {
		return nil, fmt.Errorf("only author can disable comments")
	}

	err = r.repo.DisableComments(ctx, args.PostID)
	if err != nil {
		return nil, fmt.Errorf("failed to disable comments: %w", err)
	}

	return true, nil
}
