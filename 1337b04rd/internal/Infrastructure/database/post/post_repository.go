package post

import "1337b04rd/internal/domain/models"

type PostRepository interface {
	CreatePost(post *models.Post) error
	GetPostByID(id int) (*models.Post, error)
	DeletePost(id int) error

	GetActivePosts() ([]*models.Post, error)
	GetArchivedPosts() ([]*models.Post, error)

	// UpdatePost(post *models.Post) error
	UpdateLastCommentAt(postID int) error
}


// func RegisterPostRoutes(mux *http.ServeMux, handler *post.PostHandler) {
// 	mux.HandleFunc("/posts/create", handler.CreatePostHandler)
// 	mux.HandleFunc("/posts/get/", handler.GetPostByIDHandler)
// 	mux.HandleFunc("/posts/delete/", handler.DeletePostHandler)
// 	mux.HandleFunc("/posts/active", handler.GetActivePostsHandler)
// 	mux.HandleFunc("/posts/archived", handler.GetArchivedPostsHandler)
// 	mux.HandleFunc("/posts/update-last-comment/", handler.UpdateLastCommentAtHandler)
// }
