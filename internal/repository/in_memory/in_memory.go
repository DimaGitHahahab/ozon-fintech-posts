package in_memory

import (
	"context"
	"sync"

	"github.com/DimaGitHahahab/ozon-fintech-posts/internal/domain"
	"github.com/DimaGitHahahab/ozon-fintech-posts/internal/repository"
)

// inMemoryRepository is an in-memory implementation of repository.Repository
type inMemoryRepository struct {
	mu        sync.RWMutex            // concurrent map access protection
	posts     map[int]*domain.Post    // post id -> post
	comments  map[int]*domain.Comment // comment id -> comment
	postID    int                     // autoincrement
	commentID int                     // autoincrement
}

func New() repository.Repository {
	return &inMemoryRepository{
		posts:     make(map[int]*domain.Post),
		comments:  make(map[int]*domain.Comment),
		postID:    0,
		commentID: 0,
	}
}

func (r *inMemoryRepository) GetPosts(_ context.Context) ([]*domain.Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	posts := make([]*domain.Post, 0, len(r.posts))
	for _, post := range r.posts {
		posts = append(posts, post)
	}
	return posts, nil
}

func (r *inMemoryRepository) GetPost(_ context.Context, id int) (*domain.Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	post, _ := r.posts[id]
	return post, nil
}

func (r *inMemoryRepository) ContainsPost(_ context.Context, id int) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.posts[id]
	return exists, nil
}

func (r *inMemoryRepository) CreatePost(_ context.Context, post *domain.Post) (*domain.Post, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.postID++
	post.ID = r.postID

	r.posts[post.ID] = post

	return post, nil
}

func (r *inMemoryRepository) CreateComment(_ context.Context, comment *domain.Comment) (*domain.Comment, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.commentID++
	comment.ID = r.commentID

	r.comments[comment.ID] = comment

	return comment, nil
}

func (r *inMemoryRepository) ContainsComment(_ context.Context, id int) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.comments[id]
	return exists, nil
}

func (r *inMemoryRepository) GetCommentsByPost(ctx context.Context, postID, limit, offset int) ([]*domain.Comment, error) {
	return r.getComments(func(c *domain.Comment) bool {
		return c.PostID == postID
	}, limit, offset)
}

func (r *inMemoryRepository) GetCommentsByParent(ctx context.Context, parentId, limit, offset int) ([]*domain.Comment, error) {
	return r.getComments(func(c *domain.Comment) bool {
		return c.ParentID != nil && *c.ParentID == parentId
	}, limit, offset)
}

func (r *inMemoryRepository) getComments(filterFunc func(*domain.Comment) bool, limit, offset int) ([]*domain.Comment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	comments := make([]*domain.Comment, 0, limit)
	count := 0
	for _, comment := range r.comments {
		if filterFunc(comment) {
			if count >= offset {
				comments = append(comments, comment)
				if len(comments) == limit {
					break
				}
			}
			count++
		}
	}

	return comments, nil
}

func (r *inMemoryRepository) DisableComments(_ context.Context, postID int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	post, _ := r.posts[postID]
	post.CommentsDisabled = true
	return nil
}
