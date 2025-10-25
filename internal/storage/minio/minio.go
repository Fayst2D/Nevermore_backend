package minio

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"nevermore/internal/dto"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Storage interface {
	UploadPhoto(ctx context.Context, photo dto.FileInfo) (string, error)
	UploadPdf(ctx context.Context, photo dto.FileInfo) (string, error)
	UploadPage(ctx context.Context, photo dto.FileInfo, bookId, pageNumber int) (string, error)
	DownloadFile(url string) (*dto.FileInfo, error)
}

type repo struct {
	client  *minio.Client
	BaseURL string
	photos  string // исправлена опечатка: photoes -> photos
	pages   string
	pdfs    string
}

func (r *repo) UploadPhoto(ctx context.Context, photo dto.FileInfo) (string, error) {
	ext := filepath.Ext(photo.Header.Filename)
	objectName := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	return r.upload(ctx, photo, r.photos, objectName) // исправлено: r.photoes -> r.photos
}

func (r *repo) UploadPdf(ctx context.Context, photo dto.FileInfo) (string, error) {
	ext := filepath.Ext(photo.Header.Filename)
	objectName := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	return r.upload(ctx, photo, r.pdfs, objectName)
}

func (r *repo) UploadPage(ctx context.Context, photo dto.FileInfo, bookId, pageNumber int) (string, error) {
	return r.upload(ctx, photo, r.pages, fmt.Sprintf("%d/%d", bookId, pageNumber))
}

func New(cfg Config) (Storage, error) {
	client, err := minio.New(cfg.BaseURL, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}

	fmt.Println(cfg.Pages)

	repo := &repo{
		client:  client,
		photos:  cfg.Photoes,
		pages:   cfg.Pages,
		pdfs:    cfg.Pdfs,
		BaseURL: cfg.BaseURL,
	}

	return repo, nil
}

func (r *repo) upload(ctx context.Context, photo dto.FileInfo, bucket string, path string) (string, error) {
	contentType := photo.Header.Header.Get("Content-Type")

	// Используем io.ReadAll вместо ioutil.ReadAll
	fileBytes, err := io.ReadAll(photo.File)
	if err != nil {
		return "", err
	}

	_, err = r.client.PutObject(
		ctx,
		bucket,
		path,
		bytes.NewReader(fileBytes),
		int64(len(fileBytes)),
		minio.PutObjectOptions{ContentType: contentType},
	)
	if err != nil {
		return "", fmt.Errorf("minio err: %s, %s", err, bucket)
	}

	publicURL := fmt.Sprintf("%s/%s/%s", r.BaseURL, bucket, path)
	return publicURL, nil
}

func (r *repo) DownloadFile(fileURL string) (*dto.FileInfo, error) {
	// Валидация URL
	parsedURL, err := url.Parse(fileURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %v", err)
	}

	// Разрешаем только HTTP протокол
	if parsedURL.Scheme != "" && parsedURL.Scheme != "http" {
		return nil, fmt.Errorf("only HTTP protocol is allowed")
	}

	// Если схема не указана, добавляем http://
	if parsedURL.Scheme == "" {
		// Проверяем, что это локальный URL (MinIO)
		if !strings.Contains(parsedURL.Host, "localhost") &&
			!strings.Contains(parsedURL.Host, "127.0.0.1") &&
			!strings.Contains(parsedURL.Host, "minio") &&
			!strings.HasPrefix(parsedURL.Host, "192.168.") &&
			!strings.HasPrefix(parsedURL.Host, "10.") {
			return nil, fmt.Errorf("only local network URLs are allowed")
		}
		fileURL = "http://" + fileURL
	}

	// Создаем HTTP клиент с таймаутом
	client := &http.Client{
		Timeout: 30 * time.Second, // добавляем импорт "time"
	}

	// Создаем безопасный запрос
	req, err := http.NewRequest("GET", fileURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Устанавливаем безопасные заголовки
	req.Header.Set("User-Agent", "Nevermore-App/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to download file: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download file: status code %d", resp.StatusCode)
	}

	// Ограничиваем размер файла (например, 10MB)
	maxSize := int64(10 * 1024 * 1024)
	content, err := io.ReadAll(io.LimitReader(resp.Body, maxSize))
	if err != nil {
		return nil, fmt.Errorf("failed to read file content: %v", err)
	}

	// Проверяем, не превышен ли лимит
	if int64(len(content)) == maxSize {
		// Проверяем, есть ли еще данные
		_, err := io.ReadAll(resp.Body)
		if err == nil {
			return nil, fmt.Errorf("file size exceeds limit of 10MB")
		}
	}

	// Создаем multipart.File из содержимого
	file := &dto.BytesFile{Reader: bytes.NewReader(content)}

	// Создаем FileHeader
	filename := extractFilenameFromURL(fileURL)
	fileHeader := &multipart.FileHeader{
		Filename: filename,
		Size:     int64(len(content)),
		Header:   make(map[string][]string),
	}

	contentType := http.DetectContentType(content)
	fileHeader.Header.Set("Content-Type", contentType)

	result := &dto.FileInfo{
		Header: fileHeader,
		File:   file,
	}

	return result, nil
}

// extractFilenameFromURL извлекает имя файла из URL
func extractFilenameFromURL(url string) string {
	parts := strings.Split(url, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return "downloaded_file"
}
