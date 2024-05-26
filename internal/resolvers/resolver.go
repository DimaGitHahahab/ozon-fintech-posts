package resolvers

import (
	"context"
	"fmt"
	"time"

	"github.com/DimaGitHahahab/ozon-fintech-posts/internal/domain"
	"github.com/DimaGitHahahab/ozon-fintech-posts/internal/repository"
)

var (
	ErrPostNotFound     = fmt.Errorf("post not found")
	ErrCommentNotFound  = fmt.Errorf("comment not found")
	ErrCommentsDisabled = fmt.Errorf("comments are disabled")
	ErrNotAuthor        = fmt.Errorf("only author can disable comments")
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
	if err := r.postExists(ctx, args.ID); err != nil {
		return nil, err
	}

	post, err := r.repo.GetPost(ctx, args.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get post: %w", err)
	}

	return post, nil
}

func (r *Resolver) GetCommentsByPost(ctx context.Context, args GetCommentsArgs) (any, error) {
	if err := r.postExists(ctx, args.PostID); err != nil {
		return nil, err
	}

	comments, err := r.repo.GetCommentsByPost(ctx, args.PostID, args.Limit, args.Offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get comments: %w", err)
	}

	return comments, nil
}

func (r *Resolver) GetCommentsByParent(ctx context.Context, args GetCommentsArgs) (any, error) {
	if err := r.commentExists(ctx, args.ParentID); err != nil {
		return nil, err
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
	if err := validateComment(args.Content); err != nil {
		return nil, err
	}

	comment := &domain.Comment{
		PostID:    args.PostID,
		ParentID:  &args.ParentID,
		AuthorID:  args.AuthorID,
		Content:   args.Content,
		CreatedAt: time.Now(),
	}

	if err := r.postExists(ctx, args.PostID); err != nil {
		return nil, err
	}
	if err := r.commentsDisabled(ctx, args.PostID); err != nil {
		return nil, err
	}

	savedComment, err := r.repo.CreateComment(ctx, comment)
	if err != nil {
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}

	return savedComment, nil
}

func (r *Resolver) DisableComments(ctx context.Context, args DisableCommentsArgs) (any, error) {
	if err := r.postExists(ctx, args.PostID); err != nil {
		return nil, err
	}

	if err := r.authorizeAuthor(ctx, args.PostID, args.AuthorId); err != nil {
		return nil, err
	}

	err := r.repo.DisableComments(ctx, args.PostID)
	if err != nil {
		return nil, fmt.Errorf("failed to disable comments: %w", err)
	}

	return true, nil
}

func (r *Resolver) postExists(ctx context.Context, postID int) error {
	ok, err := r.repo.ContainsPost(ctx, postID)
	if err != nil {
		return fmt.Errorf("failed to check post existence: %w", err)
	}
	if !ok {
		return fmt.Errorf("%w: %d", ErrPostNotFound, postID)
	}
	return nil
}

func (r *Resolver) commentExists(ctx context.Context, commentID int) error {
	ok, err := r.repo.ContainsComment(ctx, commentID)
	if err != nil {
		return fmt.Errorf("failed to check comment existence: %w", err)
	}
	if !ok {
		return fmt.Errorf("%w: %d", ErrCommentNotFound, commentID)
	}
	return nil
}

func (r *Resolver) commentsDisabled(ctx context.Context, postID int) error {
	post, err := r.repo.GetPost(ctx, postID)
	if err != nil {
		return fmt.Errorf("failed to get post: %w", err)
	}

	if post.CommentsDisabled {
		return fmt.Errorf("%w: %d", ErrCommentsDisabled, postID)
	}

	return nil
}

func (r *Resolver) authorizeAuthor(ctx context.Context, postID, authorID int) error {
	post, err := r.repo.GetPost(ctx, postID)
	if err != nil {
		return fmt.Errorf("failed to get post: %w", err)
	}
	if post.AuthorID != authorID {
		return ErrNotAuthor
	}

	return nil
}
