package core

import (
	"1337b04rd/internal/Infrastructure/database/post"
	"1337b04rd/internal/domain/models"
	"1337b04rd/internal/domain/ports"
	"errors"
	"fmt"
	"strconv"
)

type PostServiceImpl struct {
	repo post.PostRepository
}

func NewPostService(repo post.PostRepository) ports.PostService {
	return &PostServiceImpl{repo: repo}
}

func (s *PostServiceImpl) CreatePostService(post *models.Post) error {
	if post.CreatedAt.IsZero() {
		return errors.New("missing post creation time")
	}
	if post.Title == "" {
		return errors.New("title is required")
	}
	if post.Content == "" {
		return errors.New("content is required")
	}
	if post.Name == "" {
		return errors.New("name is required")
	}
	if post.Avatar == "" {
		return errors.New("avatar is required")
	}

	return s.repo.CreatePost(post)
}

func (s *PostServiceImpl) GetPostByIDService(id string) (*models.Post, error) {
	postID, err := strconv.Atoi(id)
	if err != nil {
		return &models.Post{}, fmt.Errorf("invalid post ID format: %v", err)
	}
	return s.repo.GetPostByID(postID)
}

func (s *PostServiceImpl) DeletePostService(id string) error {
	postID, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("invalid post ID format: %v", err)
	}
	return s.repo.DeletePost(postID)
}

func (s *PostServiceImpl) GetActivePostsService() ([]*models.Post, error) {
	return s.repo.GetActivePosts()
}

func (s *PostServiceImpl) GetArchivedPostsService() ([]*models.Post, error) {
	return s.repo.GetArchivedPosts()
}

func (s *PostServiceImpl) UpdatePostService(post *models.Post) error {
	return s.repo.UpdatePost(post)
}

func (s *PostServiceImpl) UpdateLastCommentAtService(id string) error {
	postID, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("invalid post ID format: %v", err)
	}
	return s.repo.UpdateLastCommentAt(postID)
}
