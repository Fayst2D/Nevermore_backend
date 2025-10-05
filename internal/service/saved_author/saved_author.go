package saved_author

import (
	"context"
	"fmt"
	"nevermore/internal/dto"
	model "nevermore/internal/model/saved_author"
	"nevermore/internal/storage"
	"time"
)

type Service interface {
	Create(ctx context.Context, request dto.SaveAuthorRequest) error
	Delete(ctx context.Context, authorId int) error
	GetSavedAuthorsList(ctx context.Context, userID int) ([]*dto.AuthorGetResponse, error)
}

type service struct {
	st storage.Storage
}

func New(st storage.Storage) Service {
	result := &service{
		st: st,
	}

	return result
}

func (s *service) Create(ctx context.Context, request dto.SaveAuthorRequest) error {
	author := model.SavedAuthor{
		UserId:    request.UserId,
		AuthorId:  request.AuthorId,
		CreatedAt: time.Now().UTC(),
	}

	return s.st.DB().SavedAuthor().Create(ctx, &author)
}

func (s *service) Delete(ctx context.Context, authorId int) error {
	err := s.st.DB().Author().Delete(ctx, authorId)
	if err != nil {
		return fmt.Errorf("AuthorService:Delete err -> %s", err.Error())
	}

	return nil
}

func (s *service) GetSavedAuthorsList(ctx context.Context, userID int) ([]*dto.AuthorGetResponse, error) {
	result, err := s.st.DB().SavedAuthor().GetSavedAuthorsList(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("AuthorService:GetSavedAuthorsList err -> %s", err.Error())
	}

	return result, nil
}
