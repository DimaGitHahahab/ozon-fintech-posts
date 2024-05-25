package queries

import (
	"context"
	"fmt"

	"github.com/DimaGitHahahab/ozon-fintech-posts/internal/domain"
)

const selectPosts = `
SELECT id, title, content, author_id, created_at, comments_disabled
FROM posts;
`

func (q *Queries) GetPosts(ctx context.Context) ([]*domain.Post, error) {
	rows, err := q.pool.Query(ctx, selectPosts)
	if err != nil {
		return nil, fmt.Errorf("can't select all posts: %w", err)
	}
	defer rows.Close()

	posts := make([]*domain.Post, 0)
	for rows.Next() {
		var post domain.Post
		err = rows.Scan(&post.ID, &post.Title, &post.Content, &post.AuthorID, &post.CreatedAt, &post.CommentsDisabled)
		if err != nil {
			return nil, fmt.Errorf("can't scan post row: %w", err)
		}
		posts = append(posts, &post)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error while reading rows: %w", err)
	}

	return posts, nil
}

const selectPost = `
SELECT id, title, content, author_id, created_at, comments_disabled
FROM posts
WHERE id = $1;
`

func (q *Queries) GetPost(ctx context.Context, id int) (*domain.Post, error) {
	row := q.pool.QueryRow(ctx, selectPost, id)

	var post domain.Post
	err := row.Scan(&post.ID, &post.Title, &post.Content, &post.AuthorID, &post.CreatedAt, &post.CommentsDisabled)
	if err != nil {
		return nil, fmt.Errorf("can't scan post row: %w", err)
	}

	return &post, nil
}

const insertPost = `
INSERT INTO posts
(title, content, author_id, created_at, comments_disabled)
VALUES ($1, $2, $3, $4, $5)
RETURNING id
`

func (q *Queries) CreatePost(ctx context.Context, post *domain.Post) (*domain.Post, error) {
	row := q.pool.QueryRow(ctx, insertPost, post.Title, post.Content, post.AuthorID, post.CreatedAt, post.CommentsDisabled)
	if err := row.Scan(&post.ID); err != nil {
		return nil, fmt.Errorf("can't scan post id: %w", err)
	}

	return post, nil
}

const updateDisableComments = `
UPDATE posts
SET comments_disabled = FALSE
WHERE id = $1
`

func (q *Queries) DisableComments(ctx context.Context, postID int) error {
	if _, err := q.pool.Exec(ctx, updateDisableComments, postID); err != nil {
		return fmt.Errorf("can't update row to disable comments: %w", err)
	}

	return nil
}
