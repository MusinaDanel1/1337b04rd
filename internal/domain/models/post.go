package models

import "time"

type Post struct {
	ID            int        `json:"id"`
	Title         string     `json:"title"`
	Content       string     `json:"content"`
	Avatar        string     `json:"avatar"`
	Name          string     `json:"name"`
	Image         *string    `json:"image,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	LastCommentAt *time.Time `json:"last_comment_at,omitempty"`
}
