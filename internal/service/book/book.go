package book

import (
	"bytes"
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"image/jpeg"
	"io"
	"log"
	"mime/multipart"
	"nevermore/internal/dto"
	"nevermore/internal/storage"
	"os"
	"time"

	"github.com/gammazero/workerpool"
	"github.com/gen2brain/go-fitz"
)

type Service interface {
	Create(ctx context.Context, req *dto.CreateBookRequest, file dto.FileInfo) error
}

type service struct {
	st storage.Storage
	wp *workerpool.WorkerPool
}

func New(st storage.Storage, wp *workerpool.WorkerPool) Service {
	result := &service{
		st: st,
		wp: wp,
	}

	return result
}

func (s *service) Create(ctx context.Context, req *dto.CreateBookRequest, file dto.FileInfo) error {
	tx, err := s.st.DB().BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("BookService:Create err -> %s", err.Error())
	}
	defer tx.Rollback()

	req.FileUrl, err = s.st.Cloud().UploadPdf(ctx, file)
	if err != nil {
		return fmt.Errorf("BookService:Create err -> %s", err.Error())
	}

	id, err := s.st.DB().Book().Create(ctx, tx, req)
	if err != nil {
		return fmt.Errorf("BookService:Create err -> %s", err.Error())
	}

	bookId := id
	fileUrl := req.FileUrl

	s.wp.Submit(func() {
		processCtx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
		defer cancel()

		if err := s.processBook(processCtx, tx, fileUrl, bookId); err != nil {
			log.Printf("Error processing book %d: %v", bookId, err)
		}
	})

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("BookService:Create err -> %s", err.Error())
	}

	return err
}

func (s *service) splitPdfToPages(ctx context.Context, pdfFile *dto.FileInfo, bookId int) ([]string, error) {
	tempFile, err := os.CreateTemp("", "pdf-*.pdf")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	if _, err := io.Copy(tempFile, pdfFile.File); err != nil {
		return nil, fmt.Errorf("failed to copy PDF content: %v", err)
	}

	doc, err := fitz.New(tempFile.Name())
	if err != nil {
		return nil, fmt.Errorf("failed to open PDF: %v", err)
	}
	defer doc.Close()

	totalPages := doc.NumPage()
	var pageURLs []string

	for n := 0; n < totalPages; n++ {
		img, err := doc.Image(n)
		if err != nil {
			return nil, fmt.Errorf("failed to render page %d: %v", n, err)
		}

		var imgBuf bytes.Buffer
		err = jpeg.Encode(&imgBuf, img, &jpeg.Options{Quality: 90})
		if err != nil {
			return nil, fmt.Errorf("failed to encode page %d: %v", n, err)
		}

		data := imgBuf.Bytes()
		file := &dto.BytesFile{Reader: bytes.NewReader(data)}

		header := &multipart.FileHeader{
			Filename: fmt.Sprintf("page_%d.jpg", n+1),
			Size:     int64(len(data)),
			Header:   make(map[string][]string),
		}
		header.Header.Set("Content-Type", "image/jpeg")

		pageFileInfo := dto.FileInfo{
			File:   file,
			Header: header,
		}

		pageURL, err := s.st.Cloud().UploadPage(ctx, pageFileInfo, bookId, n+1)
		if err != nil {
			return nil, fmt.Errorf("failed to upload page %d: %v", n+1, err)
		}

		pageURLs = append(pageURLs, pageURL)
	}

	return pageURLs, nil
}

func (s *service) processBook(ctx context.Context, tx *sqlx.Tx, url string, bookId int) error {
	fileInfo, err := s.st.Cloud().DownloadFile(url)
	if err != nil {
		return fmt.Errorf("BookService:Process err -> %v", err.Error())
	}

	urls, err := s.splitPdfToPages(ctx, fileInfo, bookId)

	err = s.st.DB().Book().SaveFirstPage(ctx, tx, urls[0], bookId)
	if err != nil {
		return fmt.Errorf("BookService:Process err -> %v", err.Error())
	}

	return nil
}
