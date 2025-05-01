package ports

import "1337b04rd/internal/domain/models"

type CommentService interface {
	CreateCommentService(comment *models.Comment) error
	GetCommentByIDService(id string) (*models.Comment, error)
	GetCommentsByPostIDService(postID string) ([]*models.Comment, error)
	DeleteCommentService(id string) error
	GetRepliesService(commentID string) ([]*models.Comment, error)
}
