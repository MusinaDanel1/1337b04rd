package post

import (
	"1337b04rd/internal/ui/http/post"
	"net/http"
)

// SetupRoutes initializes all routes for the PostHandler
func RegisterPostRoutes(mux *http.ServeMux, handler *post.PostHandler) {
	mux.HandleFunc("/posts/create", handler.CreatePostHandler)
	mux.HandleFunc("/posts/get/{id}", handler.GetPostByIDHandler)
}
