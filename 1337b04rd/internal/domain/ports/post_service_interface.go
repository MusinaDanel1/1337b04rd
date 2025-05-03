package ports

import "1337b04rd/internal/domain/models"

type PostService interface {
	CreatePostService(post *models.Post) error
	GetPostByIDService(id string) (*models.Post, error)
	DeletePostService(id string) error

	GetActivePostsService() ([]*models.Post, error)
	GetArchivedPostsService() ([]*models.Post, error)

	// UpdatePostService(post *models.Post) error
	UpdateLastCommentAtService(postID string) error
}
