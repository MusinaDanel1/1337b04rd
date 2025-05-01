package core

import (
	"1337b04rd/internal/Infrastructure/database/comment"
	"1337b04rd/internal/domain/models"
	"1337b04rd/internal/domain/ports"
)

type PostgresCommentService struct {
	commentRepo comment.CommentRepository
}

func NewPostgresCommentService(commentRepo comment.CommentRepository) ports.CommentService {
	return &PostgresCommentService{commentRepo: commentRepo}
}

func (s *PostgresCommentService) CreateCommentService(comment *models.Comment) error {
	return s.commentRepo.CreateComment(comment)
}

func (s *PostgresCommentService) GetCommentByIDService(id string) (*models.Comment, error) {
	return s.commentRepo.GetCommentByID(id)
}

func (s *PostgresCommentService) GetCommentsByPostIDService(postID string) ([]*models.Comment, error) {
	return s.commentRepo.GetCommentsByPostID(postID)
}

func (s *PostgresCommentService) DeleteCommentService(id string) error {
	return s.commentRepo.DeleteComment(id)
}

func (s *PostgresCommentService) GetRepliesService(commentID string) ([]*models.Comment, error) {
	return s.commentRepo.GetReplies(commentID)
}
