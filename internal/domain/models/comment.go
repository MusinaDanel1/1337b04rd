package models

import (
	"time"
)

type Comment struct {
	ID        string    `json:"id"`
	PostID    string    `json:"post_id"`
	ParentID  *string   `json:"parent_id,omitempty"`
	Content   string    `json:"content"`
	Avatar    string    `json:"avatar"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}
