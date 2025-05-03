package comment

import "1337b04rd/internal/domain/models"

type CommentRepository interface {
	CreateComment(comment *models.Comment) error
	GetCommentByID(id int) (*models.Comment, error)
	GetCommentsByPostID(postID int) ([]*models.Comment, error)
	DeleteComment(id int) error
	GetReplies(id int) ([]*models.Comment, error)
}
