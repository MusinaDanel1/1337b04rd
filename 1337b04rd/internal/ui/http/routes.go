package post

import (
	"1337b04rd/internal/ui/http/post"
	"net/http"
)

func RegisterPostRoutes(mux *http.ServeMux, handler *post.PostHandler) {
	mux.HandleFunc("/submit-post", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		handler.CreatePostHandler(w, r)
	})

	mux.HandleFunc("/post/{id}", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		handler.GetPostByIDHandler(w, r)
	})

	mux.HandleFunc("/posts/active", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		handler.GetActivePostsHandler(w, r)
	})

	mux.HandleFunc("/posts/archived", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		handler.GetArchivedPostsHandler(w, r)
	})

	mux.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		if 
	})

	// mux.HandleFunc("/posts/delete/", func(w http.ResponseWriter, r *http.Request) {
	// 	if r.Method != http.MethodDelete {
	// 		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	// 		return
	// 	}
	// 	handler.DeletePostHandler(w, r)
	// })

	// mux.HandleFunc("/posts/update-last-comment/", func(w http.ResponseWriter, r *http.Request) {
	// 	if r.Method != http.MethodPut {
	// 		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	// 		return
	// 	}
	// 	handler.UpdateLastCommentAtHandler(w, r)
	// })
}

func RegisterCommentRoutes (mux *http.ServeMux, handler *comment.CommentHandler) {
	mux.HandleFunc("/post/submit-comment", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		handler.CreateCommentHandler(w, r)
	})
}
