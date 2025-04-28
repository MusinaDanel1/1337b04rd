package storage

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

type StorageService interface {
	UploadImage(file io.Reader, bucket string, objectKey string) (string, error)
	GetImage(bucket string, objectKey string) (io.Reader, error)
	DeleteImage(imageID string, bucket string) error
}

type TripleSAdapter struct {
	baseURL    string
	httpClient *http.Client
}

func NewTripleSAdapter(baseURL string) *TripleSAdapter {
	return &TripleSAdapter{
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (a *TripleSAdapter) UploadImage(file io.Reader, bucket string, objectKey string) (string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", objectKey)
	if err != nil {
		return "", fmt.Errorf("error creating file: %w", err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return "", fmt.Errorf("error copy: %w", err)
	}

	err = writer.Close()
	if err != nil {
		return "", fmt.Errorf("error closing : %w", err)
	}

	url := fmt.Sprintf("%s/%s/%s", a.baseURL, bucket, objectKey)

	req, err := http.NewRequest("PUT", url, body)
	if err != nil {
		return "", fmt.Errorf("error request creating: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Content-Length", fmt.Sprintf("%d", body.Len()))

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("request sending error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("file upload error, status code: %s", resp.Status)
	}

	fileURL := fmt.Sprintf("%s/%s/%s", a.baseURL, bucket, objectKey)
	return fileURL, nil
}

func (a *TripleSAdapter) GetImage(bucket string, objectKey string) (io.Reader, error) {
	url := fmt.Sprintf("%s/%s/%s", a.baseURL, bucket, objectKey)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("GET request creating error: %w", err)
	}

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request sending error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error receiving the object: %w", err)
	}

	return resp.Body, nil
}
