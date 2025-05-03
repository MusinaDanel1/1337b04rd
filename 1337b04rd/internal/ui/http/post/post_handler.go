package post

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"1337b04rd/internal/domain/core"
	"1337b04rd/internal/domain/models"
	"1337b04rd/internal/domain/ports"
)

type PostHandler struct {
	postService  ports.PostService
	imageService core.ImageService
	commentService ports.CommentService
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

		uploadedURL, err := h.imageService.UploadAndGetURL(title, fileReader)
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
		log.Println("Error: Post ID is required")
		http.Error(w, "Post ID is required", http.StatusBadRequest)
		return
	}
	log.Printf("Received request to fetch post with ID: %s", postID)

	// Получаем пост из сервиса
	post, err := h.postService.GetPostByIDService(postID)
	if err != nil {
		log.Printf("Error fetching post with ID %s: %v", postID, err)
		if err.Error() == "post not found" {
			http.Error(w, "Post not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to fetch post", http.StatusInternalServerError)
		}
		return
	}

	comments, err := h.CommentService.GetCommentsByPostID(postID){
		if err != nil {
			log.Printf("Error fetching comments for post %s: %v", postID, err)
			http.Error(w, "Failed to fetch comments", http.StatusInternalServerError)
		    return
		}
	}

	// Создаем структуру PostWithImage для отправки
	item := models.PostWithImage{
		ID:        post.ID,
		Title:     post.Title,
		Content:   post.Content,
		Avatar:    post.Avatar,
		Name:      post.Name,
		CreatedAt: post.CreatedAt,
		Comments:  comments,
	}

	// Если у поста есть изображение, добавляем его в формате base64
	if post.Image != nil {
		// Получаем изображение через ImageService и кодируем его в base64
		encodedImage, err := h.imageService.ProcessImage(*post.Image)
		if err != nil {
			log.Printf("Error processing image for post %s: %v", postID, err)
			http.Error(w, "Failed to process image", http.StatusInternalServerError)
			return
		}

		item.ImageData = encodedImage
	} else {
		log.Printf("No image associated with post %s", postID)
	}

	// Устанавливаем заголовок Content-Type для отправки JSON
	w.Header().Set("Content-Type", "application/json")

	// Отправляем данные в формате JSON
	err = json.NewEncoder(w).Encode(item)
	if err != nil {
		log.Printf("Error encoding post data for post %s: %v", postID, err)
		http.Error(w, "Failed to encode post", http.StatusInternalServerError)
	}
}

func (h *PostHandler) GetActivePostsHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем активные посты через сервис
	posts, err := h.postService.GetActivePostsService()
	if err != nil {
		http.Error(w, "Failed to fetch active posts", http.StatusInternalServerError)
		return
	}

	var result []models.PostWithImage

	// Перебираем все посты
	for _, post := range posts {
		item := models.PostWithImage{
			ID:        post.ID,
			Title:     post.Title,
			Content:   post.Content,
			Avatar:    post.Avatar,
			Name:      post.Name,
			CreatedAt: post.CreatedAt,
		}

		// Если у поста есть изображение, добавляем его в формате base64
		if post.Image != nil {
			// Получаем изображение через ImageService
			encodedImage, err := h.imageService.ProcessImage(*post.Image)
			if err != nil {
				// Логируем ошибку и отправляем ответ с ошибкой
				log.Printf("Error processing image for post %d: %v", post.ID, err)
				http.Error(w, "Failed to process image", http.StatusInternalServerError)
				return
			}

			// Добавляем закодированное изображение в структуру
			item.ImageData = encodedImage
		}

		// Добавляем пост с изображением в результат
		result = append(result, item)
	}

	// Устанавливаем заголовок Content-Type для отправки JSON
	w.Header().Set("Content-Type", "application/json")

	// Отправляем результат в формате JSON
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		log.Printf("Error encoding active posts: %v", err)
		http.Error(w, "Failed to encode posts", http.StatusInternalServerError)
	}
}

func (h *PostHandler) GetArchivedPostsHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем архивные посты через сервис
	posts, err := h.postService.GetArchivedPostsService()
	if err != nil {
		http.Error(w, "Failed to fetch archived posts", http.StatusInternalServerError)
		return
	}

	var result []models.PostWithImage

	// Перебираем все посты
	for _, post := range posts {
		item := models.PostWithImage{
			ID:        post.ID,
			Title:     post.Title,
			Content:   post.Content,
			Avatar:    post.Avatar,
			Name:      post.Name,
			CreatedAt: post.CreatedAt,
		}

		// Если у поста есть изображение — обрабатываем его
		if post.Image != nil {
			encodedImage, err := h.imageService.ProcessImage(*post.Image)
			if err != nil {
				log.Printf("Error processing image for post %d: %v", post.ID, err)
				http.Error(w, "Failed to process image", http.StatusInternalServerError)
				return
			}
			item.ImageData = encodedImage
		}

		result = append(result, item)
	}

	// Устанавливаем заголовок и отправляем JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		log.Printf("Error encoding archived posts: %v", err)
		http.Error(w, "Failed to encode posts", http.StatusInternalServerError)
	}
}

// func (h *PostHandler) DeletePostHandler(w http.ResponseWriter, r *http.Request) {
// 	parts := strings.Split(r.URL.Path, "/")
// 	postID := parts[len(parts)-1]

// 	if postID == "" {
// 		http.Error(w, "Post ID is required", http.StatusBadRequest)
// 		return
// 	}

// 	err := h.postService.DeletePostService(postID)
// 	if err != nil {
// 		if err.Error() == "post not found" {
// 			http.Error(w, "Post not found", http.StatusNotFound)
// 		} else {
// 			http.Error(w, "Failed to delete post", http.StatusInternalServerError)
// 		}
// 		return
// 	}

// 	w.WriteHeader(http.StatusNoContent)
// }
