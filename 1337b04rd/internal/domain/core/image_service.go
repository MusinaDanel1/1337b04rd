package core

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"1337b04rd/internal/Infrastructure/storage"
)

type ImageService struct {
	storageAdapter storage.StorageService
}

func NewImageService(storageAdapter storage.StorageService) *ImageService {
	return &ImageService{storageAdapter: storageAdapter}
}

// UploadAndGetURL загружает изображение в бакет и возвращает публичный URL
func (s *ImageService) UploadAndGetURL(title string, file io.Reader) (string, error) {
	bucketName := "posts"
	currentTime := time.Now()
	objectKey := fmt.Sprintf("%s-%d.jpg", strings.ReplaceAll(title, " ", "_"), currentTime.UnixNano())

	log.Printf("[ImageService] Загрузка изображения в бакет: %s, объект: %s", bucketName, objectKey)

	uploadedURL, err := s.storageAdapter.UploadImage(file, bucketName, objectKey)
	if err != nil {
		log.Printf("[ImageService] Ошибка при загрузке изображения: %v", err)
		return "", fmt.Errorf("upload failed: %w", err)
	}

	log.Printf("[ImageService] Изображение успешно загружено по URL: %s", uploadedURL)
	return uploadedURL, nil
}

// ProcessImage получает изображение из хранилища и кодирует его в base64
func (s *ImageService) ProcessImage(imageURL string) (*string, error) {
	log.Printf("[ImageService] Начало обработки изображения: %s", imageURL)

	// Убираем префикс http:// или https://
	if strings.HasPrefix(imageURL, "http://") {
		imageURL = imageURL[len("http://"):]
		log.Printf("[ImageService] Удален префикс http:// -> %s", imageURL)
	} else if strings.HasPrefix(imageURL, "https://") {
		imageURL = imageURL[len("https://"):]
		log.Printf("[ImageService] Удален префикс https:// -> %s", imageURL)
	}

	// Разбиваем строку на части
	parts := strings.SplitN(imageURL, "/", 3)
	log.Printf("[ImageService] Части после split: %v", parts)

	// Проверяем, что части пути корректные
	if len(parts) != 3 {
		errMsg := "[ImageService] Неверный формат URL изображения"
		log.Println(errMsg)
		return nil, fmt.Errorf("invalid image URL format")
	}

	bucket := parts[1]
	objectKey := parts[2]

	log.Printf("[ImageService] Получение изображения из бакета '%s' с ключом '%s'", bucket, objectKey)

	// Создаем адаптер для работы с хранилищем
	storageAdapter := storage.NewTripleSAdapter("http://triple-s:9000")
	reader, err := storageAdapter.GetImage(bucket, objectKey)
	if err != nil {
		log.Printf("[ImageService] Ошибка получения изображения: %v", err)
		return nil, fmt.Errorf("failed to get image: %w", err)
	}
	defer reader.Close()

	// Читаем данные изображения
	imageData, err := io.ReadAll(reader)
	if err != nil {
		log.Printf("[ImageService] Ошибка чтения данных изображения: %v", err)
		return nil, fmt.Errorf("failed to read image data: %w", err)
	}

	// Кодируем изображение в base64
	encodedImage := base64.StdEncoding.EncodeToString(imageData)
	log.Printf("[ImageService] Изображение успешно закодировано в base64 (длина: %d символов)", len(encodedImage))

	return &encodedImage, nil
}
