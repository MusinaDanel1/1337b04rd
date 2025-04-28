package post

import (
	"1337b04rd/internal/Infrastructure/storage"
	"1337b04rd/internal/domain/models"
	"1337b04rd/internal/domain/ports"
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"time"
)

type PostHandler struct {
	postService ports.PostService
}

func NewPostHandler(service ports.PostService) *PostHandler {
	return &PostHandler{postService: service}
}

func (h *PostHandler) CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Failed to parse form: "+err.Error(), http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
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

		storageAdapter := storage.NewTripleSAdapter("http://your-triple-s-service-url")
		bucketName := "posts"
		objectKey := "image.jpg"

		uploadedURL, err := storageAdapter.UploadImage(fileReader, bucketName, objectKey)
		if err != nil {
			http.Error(w, "Failed to upload image: "+err.Error(), http.StatusInternalServerError)
			return
		}
		imageURL = &uploadedURL
	}

	post := &models.Post{
		Title:     title,
		Content:   content,
		Avatar:    "https://rickandmortyapi.com/api/character/avatar/1.jpeg",
		Name:      "Rick Sanchez",
		Image:     imageURL,
		CreatedAt: time.Now(),
	}

	if err := h.postService.CreatePostService(post); err != nil {
		http.Error(w, "Failed to save post: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// GetPostByID handles HTTP GET requests to retrieve a post by its ID.
// Expects an `id` parameter in the query string.
func (h *PostHandler) GetPostByIDHandler(w http.ResponseWriter, r *http.Request) {
	postID := r.PathValue("id")
	if postID == "" {
		http.Error(w, "Post ID is required", http.StatusBadRequest)
		return
	}

	post, err := h.postService.GetPostByIDService(postID)
	if err != nil {
		if err.Error() == "post not found" {
			http.Error(w, "Post not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to fetch post", http.StatusInternalServerError)
		}
		return
	}

	var imageData []byte
	if post.Image != nil {
		parts := strings.SplitN(*post.Image, "/", 2)
		if len(parts) != 2 {
			http.Error(w, "Invalid image path format", http.StatusInternalServerError)
			return
		}
		bucket := parts[0]
		key := parts[1]

		storageAdapter := storage.NewTripleSAdapter("http://your-triple-s-service-url")
		reader, err := storageAdapter.GetImage(bucket, key)
		if err != nil {
			http.Error(w, "Failed to download image from storage", http.StatusInternalServerError)
			return
		}
		imageData, err = io.ReadAll(reader)
		if err != nil {
			http.Error(w, "Failed to read image data", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "multipart/form-data")

	writer := multipart.NewWriter(w)
	defer writer.Close()

	if part, err := writer.CreateFormField("post"); err == nil {
		json.NewEncoder(part).Encode(post)
	}

	if imageData != nil {
		if part, err := writer.CreateFormFile("image", "image.jpg"); err == nil {
			part.Write(imageData)
		}
	}
}

// handleImageUpload uploads raw image data to object storage (e.g. S3 or MinIO)
// and returns the URL/path of the stored image.
//
// TODO: Replace stub with actual implementation using AWS or MinIO SDK.

// handleImageDownload downloads the image from object storage given its URL/path
// and returns the raw byte content of the image.
//
// TODO: Replace stub with actual implementation using AWS or MinIO SDK.
