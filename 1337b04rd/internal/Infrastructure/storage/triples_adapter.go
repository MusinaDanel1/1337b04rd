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

// ensureBucket creates the bucket if it does not exist
func (a *TripleSAdapter) ensureBucket(bucket string) error {
	url := fmt.Sprintf("%s/%s", a.baseURL, bucket)
	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		return fmt.Errorf("creating bucket request: %w", err)
	}
	resp, err := a.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("sending bucket request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusConflict {
		return fmt.Errorf("bucket creation failed: %s", resp.Status)
	}
	return nil
}

func (a *TripleSAdapter) UploadImage(file io.Reader, bucket string, objectKey string) (string, error) {
	if err := a.ensureBucket(bucket); err != nil {
		return "", err
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", objectKey)
	if err != nil {
		return "", fmt.Errorf("creating form file: %w", err)
	}

	if _, err = io.Copy(part, file); err != nil {
		return "", fmt.Errorf("copying file: %w", err)
	}

	if err = writer.Close(); err != nil {
		return "", fmt.Errorf("closing writer: %w", err)
	}

	url := fmt.Sprintf("%s/%s/%s", a.baseURL, bucket, objectKey)
	req, err := http.NewRequest("PUT", url, body)
	if err != nil {
		return "", fmt.Errorf("creating PUT request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("upload request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("file upload failed: %s", resp.Status)
	}

	return url, nil
}

func (a *TripleSAdapter) GetImage(bucket string, objectKey string) (io.ReadCloser, error) {
	url := fmt.Sprintf("%s/%s/%s", a.baseURL, bucket, objectKey)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating GET request: %w", err)
	}

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("sending GET request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		return nil, fmt.Errorf("object retrieval failed: %s", resp.Status)
	}

	return resp.Body, nil // ← io.ReadCloser
}
