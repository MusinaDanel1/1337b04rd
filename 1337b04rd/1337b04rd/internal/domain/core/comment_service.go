package core

import (
	"1337b04rd/internal/Infrastructure/database/comment"
	"1337b04rd/internal/domain/models"
	"1337b04rd/internal/domain/ports"
	"errors"
	"fmt"
	"strconv"
)

type PostgresCommentService struct {
	commentRepo comment.CommentRepository
}

func NewPostgresCommentService(commentRepo comment.CommentRepository) ports.CommentService {
	return &PostgresCommentService{commentRepo: commentRepo}
}

func (s *PostgresCommentService) CreateCommentService(comment *models.Comment) error {
	if comment.Content == "" {
		return errors.New("content is required")
	}
	if comment.Name == "" {
		return errors.New("name is required")
	}
	if comment.Avatar == "" {
		return errors.New("avatar is required")
	}
	return s.commentRepo.CreateComment(comment)
}

func (s *PostgresCommentService) GetCommentByIDService(id string) (*models.Comment, error) {
	commentID, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("invalid comment ID format: %v", err)
	}
	return s.commentRepo.GetCommentByID(commentID)
}

func (s *PostgresCommentService) GetCommentsByPostIDService(postID string) ([]*models.Comment, error) {
	postid, err := strconv.Atoi(postID)
	if err != nil {
		return nil, fmt.Errorf("invalid post ID format: %v", err)
	}
	return s.commentRepo.GetCommentsByPostID(postid)
}

func (s *PostgresCommentService) DeleteCommentService(id string) error {
	commentID, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("invalid comment ID format: %v", err)
	}
	return s.commentRepo.DeleteComment(commentID)
}

func (s *PostgresCommentService) GetRepliesService(id string) ([]*models.Comment, error) {
	commentID, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("invalid comment ID format: %v", err)
	}
	return s.commentRepo.GetReplies(commentID)
}
