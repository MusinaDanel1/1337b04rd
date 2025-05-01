package comment

import "1337b04rd/internal/domain/models"

type CommentRepository interface {
	CreateComment(comment *models.Comment) error
	GetCommentByID(id string) (*models.Comment, error)
	GetCommentsByPostID(postID string) ([]*models.Comment, error)
	DeleteComment(id string) error
	GetReplies(commentID string) ([]*models.Comment, error)
}
