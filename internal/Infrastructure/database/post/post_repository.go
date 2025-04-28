package post

import "1337b04rd/internal/domain/models"

type PostRepository interface {
	CreatePost(post *models.Post) error
	GetPostByID(id int) (*models.Post, error)
	DeletePost(id int) error

	GetActivePosts() ([]*models.Post, error)
	GetArchivedPosts() ([]*models.Post, error)

	UpdatePost(post *models.Post) error
	UpdateLastCommentAt(postID int) error
}
