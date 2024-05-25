package queries

import (
	"context"
	"fmt"

	"github.com/DimaGitHahahab/ozon-fintech-posts/internal/domain"
)

const insertComment = `
INSERT INTO comments
(post_id, parent_id, author_id, content, created_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING id
`

func (q *Queries) CreateComment(ctx context.Context, comment *domain.Comment) (*domain.Comment, error) {
	row := q.pool.QueryRow(ctx, insertComment,
		comment.PostID, nil, comment.AuthorID, comment.Content, comment.CreatedAt)

	if err := row.Scan(&comment.ID); err != nil {
		return nil, fmt.Errorf("can't scan comment id: %w", err)
	}

	return comment, nil
}

const selectCommentsByPost = `
SELECT id, post_id, parent_id, author_id, content, created_at
FROM comments
WHERE post_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3
`

func (q *Queries) GetCommentsByPost(ctx context.Context, postID int, limit, offset int) ([]*domain.Comment, error) {
	rows, err := q.pool.Query(ctx, selectCommentsByPost, postID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("can't select comments: %w", err)
	}
	defer rows.Close()

	var comments []*domain.Comment
	for rows.Next() {
		var c domain.Comment
		if err := rows.Scan(&c.ID, &c.PostID, &c.ParentID, &c.AuthorID, &c.Content, &c.CreatedAt); err != nil {
			return nil, fmt.Errorf("can't scan comment: %w", err)
		}
		comments = append(comments, &c)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error while reading rows: %w", err)
	}

	return comments, nil
}

const selectCommentsByParent = `
SELECT id, post_id, parent_id, author_id, content, created_at
FROM comments
WHERE parent_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3
`

func (q *Queries) GetCommentsByParent(ctx context.Context, parentId int, limit, offset int) ([]*domain.Comment, error) {
	rows, err := q.pool.Query(ctx, selectCommentsByParent, parentId, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("can't selecect comments: %w", err)
	}
	defer rows.Close()

	var comments []*domain.Comment
	for rows.Next() {
		var c domain.Comment
		if err := rows.Scan(&c.ID, &c.PostID, &c.ParentID, &c.AuthorID, &c.Content, &c.CreatedAt); err != nil {
			return nil, fmt.Errorf("can't scan comment: %w", err)
		}
		comments = append(comments, &c)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error while reading rows: %w", err)
	}

	return comments, nil
}
