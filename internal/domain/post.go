package domain

import "time"

type Post struct {
	ID               int       `json:"id"`
	Title            string    `json:"title"`
	Content          string    `json:"content"`
	AuthorID         int       `json:"author_id"`
	CreatedAt        time.Time `json:"created_at"`
	CommentsDisabled bool      `json:"comments_disabled"`
}
