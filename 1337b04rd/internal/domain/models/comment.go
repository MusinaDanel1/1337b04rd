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
	Image     *string   `json:"image,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type CommentWithImage struct {
	ID        int       `json:"id"`
	PostID    int       `json:"post_id"`
	Content   string    `json:"content"`
	Avatar    string    `json:"avatar"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	ImageData string    `json:"image_data,omitempty"` // base64
}
