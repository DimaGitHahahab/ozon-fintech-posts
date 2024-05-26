package resolvers

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/DimaGitHahahab/ozon-fintech-posts/internal/domain"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
)

func TestResolver_GetPosts(t *testing.T) {
	mockRepo := NewRepositoryMock(minimock.NewController(t))

	post := &domain.Post{
		ID:               1,
		Title:            "Title",
		Content:          "Content",
		AuthorID:         1,
		CreatedAt:        time.Now(),
		CommentsDisabled: false,
	}

	mockRepo.GetPostsMock.Expect(minimock.AnyContext).Return([]*domain.Post{post}, nil)

	resolver := NewResolver(mockRepo)

	res, err := resolver.GetPosts(context.Background())
	assert.NoError(t, err)

	posts, ok := res.([]*domain.Post)
	assert.True(t, ok)

	assert.Len(t, posts, 1)
	assert.Equal(t, post, posts[0])
}

func TestResolver_GetPost(t *testing.T) {
	mockRepo := NewRepositoryMock(minimock.NewController(t))

	post := &domain.Post{
		ID:               1,
		Title:            "Title",
		Content:          "Content",
		AuthorID:         1,
		CreatedAt:        time.Now(),
		CommentsDisabled: false,
	}

	mockRepo.ContainsPostMock.Expect(minimock.AnyContext, 1).Return(true, nil)
	mockRepo.GetPostMock.Expect(minimock.AnyContext, 1).Return(post, nil)

	resolver := NewResolver(mockRepo)

	res, err := resolver.GetPost(context.Background(), PostArgs{ID: 1})
	assert.NoError(t, err)

	p, ok := res.(*domain.Post)
	assert.True(t, ok)

	assert.Equal(t, post, p)
}

func TestResolver_GetPost_NoSuchPost(t *testing.T) {
	mockRepo := NewRepositoryMock(minimock.NewController(t))

	mockRepo.ContainsPostMock.Expect(minimock.AnyContext, 1).Return(false, nil)

	resolver := NewResolver(mockRepo)

	post, err := resolver.GetPost(context.Background(), PostArgs{ID: 1})
	assert.ErrorIs(t, err, ErrPostNotFound)
	assert.Nil(t, post)
}

func TestResolver_GetCommentsByPost(t *testing.T) {
	mockRepo := NewRepositoryMock(minimock.NewController(t))

	comment := &domain.Comment{
		ID:        1,
		PostID:    1,
		ParentID:  nil,
		AuthorID:  1,
		Content:   "Content",
		CreatedAt: time.Now(),
	}

	mockRepo.ContainsPostMock.Expect(minimock.AnyContext, 1).Return(true, nil)
	mockRepo.GetCommentsByPostMock.Expect(minimock.AnyContext, 1, 10, 0).Return([]*domain.Comment{comment}, nil)

	resolver := NewResolver(mockRepo)

	res, err := resolver.GetCommentsByPost(context.Background(), GetCommentsArgs{PostID: 1, Limit: 10, Offset: 0})
	assert.NoError(t, err)

	comments, ok := res.([]*domain.Comment)
	assert.True(t, ok)

	assert.Len(t, comments, 1)
	assert.Equal(t, comment, comments[0])
}

func TestResolver_GetCommentsByPost_NoSuchPost(t *testing.T) {
	mockRepo := NewRepositoryMock(minimock.NewController(t))

	mockRepo.ContainsPostMock.Expect(minimock.AnyContext, 1).Return(false, nil)

	resolver := NewResolver(mockRepo)

	comments, err := resolver.GetCommentsByPost(context.Background(), GetCommentsArgs{PostID: 1, Limit: 10, Offset: 0})
	assert.ErrorIs(t, err, ErrPostNotFound)
	assert.Nil(t, comments)
}

func TestResolver_GetCommentsByParent(t *testing.T) {
	mockRepo := NewRepositoryMock(minimock.NewController(t))

	comment := &domain.Comment{
		ID:        1,
		PostID:    1,
		ParentID:  nil,
		AuthorID:  1,
		Content:   "Content",
		CreatedAt: time.Now(),
	}

	mockRepo.ContainsCommentMock.Expect(minimock.AnyContext, 1).Return(true, nil)
	mockRepo.GetCommentsByParentMock.Expect(minimock.AnyContext, 1, 10, 0).Return([]*domain.Comment{comment}, nil)

	resolver := NewResolver(mockRepo)

	pointerToParent := 1
	res, err := resolver.GetCommentsByParent(context.Background(), GetCommentsArgs{PostID: 1, ParentID: &pointerToParent, Limit: 10, Offset: 0})
	assert.NoError(t, err)

	comments, ok := res.([]*domain.Comment)
	assert.True(t, ok)

	assert.Len(t, comments, 1)
	assert.Equal(t, comment, comments[0])
}

func TestResolver_GetCommentsByParent_NoSuchComment(t *testing.T) {
	mockRepo := NewRepositoryMock(minimock.NewController(t))

	mockRepo.ContainsCommentMock.Expect(minimock.AnyContext, 1).Return(false, nil)

	resolver := NewResolver(mockRepo)

	pointerToParent := 1
	comments, err := resolver.GetCommentsByParent(context.Background(), GetCommentsArgs{PostID: 1, ParentID: &pointerToParent, Limit: 10, Offset: 0})
	assert.ErrorIs(t, err, ErrCommentNotFound)
	assert.Nil(t, comments)
}

func TestResolver_CreatePost(t *testing.T) {
	mockRepo := NewRepositoryMock(minimock.NewController(t))

	post := &domain.Post{
		Title:    "Title",
		Content:  "Content",
		AuthorID: 1,
	}

	expectedPost := &domain.Post{
		ID:       1,
		Title:    post.Title,
		Content:  post.Content,
		AuthorID: post.AuthorID,
	}

	mockRepo.CreatePostMock.Set(func(ctx context.Context, p *domain.Post) (*domain.Post, error) {
		p.ID = 1
		p.CreatedAt = time.Now().UTC()
		return p, nil
	})

	resolver := NewResolver(mockRepo)

	before := time.Now().UTC()
	res, err := resolver.CreatePost(context.Background(), CreatePostArgs{
		Title:    "Title",
		Content:  "Content",
		AuthorID: 1,
	})
	after := time.Now().UTC()

	assert.NoError(t, err)

	resultPost, ok := res.(*domain.Post)
	assert.True(t, ok)

	assert.Equal(t, expectedPost.ID, resultPost.ID)
	assert.Equal(t, expectedPost.Title, resultPost.Title)
	assert.Equal(t, expectedPost.Content, resultPost.Content)
	assert.Equal(t, expectedPost.AuthorID, resultPost.AuthorID)

	assert.True(t, before.Before(resultPost.CreatedAt) || before.Equal(resultPost.CreatedAt))
	assert.True(t, after.After(resultPost.CreatedAt) || after.Equal(resultPost.CreatedAt))
}

func TestResolver_DisableComments(t *testing.T) {
	mockRepo := NewRepositoryMock(minimock.NewController(t))

	mockRepo.ContainsPostMock.Expect(minimock.AnyContext, 1).Return(true, nil)
	mockRepo.GetPostMock.Expect(minimock.AnyContext, 1).Return(&domain.Post{ID: 1, AuthorID: 1}, nil)
	mockRepo.DisableCommentsMock.Expect(minimock.AnyContext, 1).Return(nil)

	resolver := NewResolver(mockRepo)

	res, err := resolver.DisableComments(context.Background(), DisableCommentsArgs{PostID: 1, AuthorId: 1})
	assert.NoError(t, err)
	ok, cast := res.(bool)
	assert.True(t, cast)
	assert.True(t, ok)
}

func TestResolver_DisableComments_NoSuchPost(t *testing.T) {
	mockRepo := NewRepositoryMock(minimock.NewController(t))

	mockRepo.ContainsPostMock.Expect(minimock.AnyContext, 1).Return(false, nil)

	resolver := NewResolver(mockRepo)

	res, err := resolver.DisableComments(context.Background(), DisableCommentsArgs{PostID: 1, AuthorId: 1})
	assert.ErrorIs(t, err, ErrPostNotFound)
	assert.Nil(t, res)
}

func TestResolver_DisableComments_WrongAuthor(t *testing.T) {
	mockRepo := NewRepositoryMock(minimock.NewController(t))

	mockRepo.ContainsPostMock.Expect(minimock.AnyContext, 1).Return(true, nil)
	mockRepo.GetPostMock.Expect(minimock.AnyContext, 1).Return(&domain.Post{ID: 1, AuthorID: 1}, nil)

	resolver := NewResolver(mockRepo)

	res, err := resolver.DisableComments(context.Background(), DisableCommentsArgs{PostID: 1, AuthorId: 2})
	assert.ErrorIs(t, err, ErrNotAuthor)
	assert.Nil(t, res)
}

func TestResolver_CreateComment(t *testing.T) {
	mockRepo := NewRepositoryMock(minimock.NewController(t))

	comment := &domain.Comment{
		PostID:    1,
		ParentID:  nil,
		AuthorID:  1,
		Content:   "Content",
		CreatedAt: time.Now(),
	}

	expectedComment := &domain.Comment{
		ID:        1,
		PostID:    comment.PostID,
		ParentID:  comment.ParentID,
		AuthorID:  comment.AuthorID,
		Content:   comment.Content,
		CreatedAt: comment.CreatedAt,
	}
	mockRepo.ContainsPostMock.Expect(minimock.AnyContext, 1).Return(true, nil)
	mockRepo.GetPostMock.Expect(minimock.AnyContext, 1).Return(&domain.Post{ID: 1, CommentsDisabled: false}, nil)

	mockRepo.CreateCommentMock.Set(func(ctx context.Context, c *domain.Comment) (*domain.Comment, error) {
		c.ID = 1
		c.CreatedAt = time.Now().UTC()
		return c, nil
	})

	resolver := NewResolver(mockRepo)

	before := time.Now().UTC()
	res, err := resolver.CreateComment(context.Background(), CreateCommentArgs{
		PostID:   1,
		ParentID: nil,
		AuthorID: 1,
		Content:  "Content",
	})
	after := time.Now().UTC()

	assert.NoError(t, err)

	resultComment, ok := res.(*domain.Comment)
	assert.True(t, ok)

	assert.Equal(t, expectedComment.ID, resultComment.ID)
	assert.Equal(t, expectedComment.PostID, resultComment.PostID)
	assert.Equal(t, expectedComment.ParentID, resultComment.ParentID)
	assert.Equal(t, expectedComment.AuthorID, resultComment.AuthorID)
	assert.Equal(t, expectedComment.Content, resultComment.Content)

	assert.True(t, before.Before(resultComment.CreatedAt) || before.Equal(resultComment.CreatedAt))
	assert.True(t, after.After(resultComment.CreatedAt) || after.Equal(resultComment.CreatedAt))
}

func TestResolver_CreateComment_CommentsForbidden(t *testing.T) {
	mockRepo := NewRepositoryMock(minimock.NewController(t))

	mockRepo.ContainsPostMock.Expect(minimock.AnyContext, 1).Return(true, nil)
	mockRepo.GetPostMock.Expect(minimock.AnyContext, 1).Return(&domain.Post{ID: 1, CommentsDisabled: true}, nil)

	resolver := NewResolver(mockRepo)

	res, err := resolver.CreateComment(context.Background(), CreateCommentArgs{
		PostID:   1,
		ParentID: nil,
		AuthorID: 1,
		Content:  "Content",
	})
	assert.ErrorIs(t, err, ErrCommentsDisabled)
	assert.Nil(t, res)
}

func TestResolver_CreateComment_NoSuchPost(t *testing.T) {
	mockRepo := NewRepositoryMock(minimock.NewController(t))

	mockRepo.ContainsPostMock.Expect(minimock.AnyContext, 1).Return(false, nil)

	resolver := NewResolver(mockRepo)

	res, err := resolver.CreateComment(context.Background(), CreateCommentArgs{
		PostID:   1,
		ParentID: nil,
		AuthorID: 1,
		Content:  "Content",
	})
	assert.ErrorIs(t, err, ErrPostNotFound)
	assert.Nil(t, res)
}

func TestResolver_CreateComment_TooLongContent(t *testing.T) {
	mockRepo := NewRepositoryMock(minimock.NewController(t))

	resolver := NewResolver(mockRepo)

	res, err := resolver.CreateComment(context.Background(), CreateCommentArgs{
		PostID:   1,
		ParentID: nil,
		AuthorID: 1,
		Content:  strings.Repeat("a", maxLength+1),
	})
	assert.ErrorIs(t, err, ErrInvalidComment)
	assert.Nil(t, res)
}

func TestResolver_BadID(t *testing.T) {
	mockRepo := NewRepositoryMock(minimock.NewController(t))

	resolver := NewResolver(mockRepo)

	res, err := resolver.CreateComment(context.Background(), CreateCommentArgs{
		PostID:   1,
		ParentID: nil,
		AuthorID: -1,
		Content:  "Content",
	})
	assert.ErrorIs(t, err, ErrNotPositiveID)
	assert.Nil(t, res)

	res, err = resolver.CreateComment(context.Background(), CreateCommentArgs{
		PostID:   -1,
		ParentID: nil,
		AuthorID: 1,
	},
	)
	assert.ErrorIs(t, err, ErrNotPositiveID)
	assert.Nil(t, res)
}
