package post

import (
	"1337b04rd/internal/domain/models"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type mockPostService struct {
	CreatedPost *models.Post
	GetPost     *models.Post
	Err         error
}

func (m *mockPostService) CreatePostService(post *models.Post) error {
	m.CreatedPost = post
	return m.Err
}

func (m *mockPostService) GetPostByIDService(id string) (*models.Post, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	return m.GetPost, nil
}

type mockImageService struct {
	UploadURL string
	ImageData string
	Err       error
}

func (m *mockImageService) UploadAndGetURL(name string, file io.Reader) (string, error) {
	return m.UploadURL, m.Err
}

func (m *mockImageService) ProcessImage(url string) (string, error) {
	return m.ImageData, m.Err
}

type mockStorage struct{}

func (m *mockStorage) Save(key string, data []byte) error {
	return nil
}

func withFakeSession(r *http.Request) *http.Request {
	ctx := context.WithValue(r.Context(), middleware.SessionContextKey, &models.Session{
		Name:   "John Doe",
		Avatar: "avatar-url",
	})
	return r.WithContext(ctx)
}

func TestCreatePostHandler_Success(t *testing.T) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	_ = writer.WriteField("title", "Test Title")
	_ = writer.WriteField("content", "Test Content")
	part, _ := writer.CreateFormFile("file", "image.jpg")
	part.Write([]byte("fake-image-data"))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/posts", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req = withFakeSession(req)

	rec := httptest.NewRecorder()

	handler := http.NewPostHandler(
		&mockPostService{},
		&mockImageService{UploadURL: "http://image.url"},
		&mockStorage{},
	)
	handler.CreatePostHandler(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", rec.Code)
	}
}

func TestCreatePostHandler_NoSession(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/posts", nil)
	rec := httptest.NewRecorder()

	handler := http.NewPostHandler(&mockPostService{}, &mockImageService{}, &mockStorage{})
	handler.CreatePostHandler(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rec.Code)
	}
}

func TestGetPostByIDHandler_Success(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/posts/123", nil)
	req = req.WithContext(context.WithValue(req.Context(), middleware.SessionContextKey, &models.Session{Name: "John"}))
	req.SetPathValue("id", "123")

	rec := httptest.NewRecorder()

	handler := http.NewPostHandler(
		&mockPostService{
			GetPost: &models.Post{
				ID:        123,
				Title:     "Test",
				Content:   "Content",
				Avatar:    "avatar",
				Name:      "John",
				Image:     ptr("http://img"),
				CreatedAt: time.Now(),
			},
		},
		&mockImageService{ImageData: "encoded-img"},
		&mockStorage{},
	)

	handler.GetPostByIDHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var response models.PostWithImage
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if response.Title != "Test" || response.ImageData != "encoded-img" {
		t.Errorf("unexpected response: %+v", response)
	}
}
