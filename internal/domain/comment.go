package domain

import "time"

type Comment struct {
	ID        int       `json:"id"`
	PostID    int       `json:"post_id"`
	ParentID  *int      `json:"parent_id,omitempty"` // nil when the comment is a root comment
	AuthorID  int       `json:"author_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}
