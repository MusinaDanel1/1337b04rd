package comment

import (
	"1337b04rd/internal/domain/core"
	"1337b04rd/internal/domain/models"
	"1337b04rd/internal/domain/ports"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

type CommentHandler struct {
	commentService ports.CommentService
	imageService   core.ImageService
}

func NewCommentHandler(service ports.CommentService) *CommentHandler {
	return &CommentHandler{commentService: service}
}

func (h *CommentHandler) CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Failed to parse form: "+err.Error(), http.StatusBadRequest)
		return
	}

	postID := r.FormValue("post_id")
	content := r.FormValue("content")

	var imageURL *string
	file, _, err := r.FormFile("file")
	if err != nil && err != http.ErrMissingFile {
		http.Error(w, "Error retrieving file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err == nil {
		defer file.Close()
		fileBytes, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, "Error reading file: "+err.Error(), http.StatusInternalServerError)
			return
		}
		fileReader := bytes.NewReader(fileBytes)
		uploadedURL, err := h.imageService.UploadAndGetURL("comment-"+postID, fileReader)
		if err != nil {
			http.Error(w, "Failed to upload image: "+err.Error(), http.StatusInternalServerError)
			return
		}
		imageURL = &uploadedURL
	}

	comment := &models.Comment{
		PostID:    atoi(postID),
		Content:   content,
		Avatar:    "https://rickandmortyapi.com/api/character/avatar/1.jpeg",
		Name:      "Rick Sanchez",
		Image:     imageURL,
		CreatedAt: time.Now(),
	}

	if err := h.commentService.CreateComment(comment); err != nil {
		http.Error(w, "Failed to create comment: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// GetCommentsByPostIDHandler получает комментарии по ID поста
func (h *CommentHandler) GetCommentsByPostIDHandler(w http.ResponseWriter, r *http.Request) {
	postID := r.PathValue("post_id")
	if postID == "" {
		http.Error(w, "Post ID is required", http.StatusBadRequest)
		return
	}

	comments, err := h.commentService.GetCommentsByPostID(postID)
	if err != nil {
		http.Error(w, "Failed to get comments: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var result []models.CommentWithImage
	for _, c := range comments {
		item := models.CommentWithImage{
			ID:        c.ID,
			PostID:    c.PostID,
			Content:   c.Content,
			Avatar:    c.Avatar,
			Name:      c.Name,
			CreatedAt: c.CreatedAt,
		}
		if c.Image != nil {
			imgData, err := h.imageService.ProcessImage(*c.Image)
			if err != nil {
				log.Printf("Failed to process image for comment %d: %v", c.ID, err)
				http.Error(w, "Failed to process image", http.StatusInternalServerError)
				return
			}
			item.ImageData = imgData
		}
		result = append(result, item)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "Failed to encode comments", http.StatusInternalServerError)
	}
}

func atoi(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}
