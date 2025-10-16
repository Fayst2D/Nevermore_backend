package minio

import (
	"bytes"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"nevermore/internal/dto"
	"path/filepath"
	"strings"
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
	photoes string
	pages   string
	pdfs    string
}

func (r *repo) UploadPhoto(ctx context.Context, photo dto.FileInfo) (string, error) {
	ext := filepath.Ext(photo.Header.Filename)
	objectName := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	return r.upload(ctx, photo, r.photoes, objectName)
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
		photoes: cfg.Photoes,
		pages:   cfg.Pages,
		pdfs:    cfg.Pdfs,
		BaseURL: cfg.BaseURL,
	}

	return repo, nil
}

func (r *repo) upload(ctx context.Context, photo dto.FileInfo, bucket string, path string) (string, error) {
	contentType := photo.Header.Header.Get("Content-Type")

	fileBytes, err := ioutil.ReadAll(photo.File)
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

func (r *repo) DownloadFile(url string) (*dto.FileInfo, error) {
	reqUrl := fmt.Sprintf("http://%s", url)

	resp, err := http.Get(reqUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to download file: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download file: status code %d", resp.StatusCode)
	}

	// Читаем содержимое файла
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read file content: %v", err)
	}

	// Создаем multipart.File из содержимого
	file := &dto.BytesFile{Reader: bytes.NewReader(content)}

	// Создаем FileHeader
	filename := extractFilenameFromURL(url)
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
