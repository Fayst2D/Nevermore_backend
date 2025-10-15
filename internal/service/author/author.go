package author

import (
	"context"
	"fmt"
	"nevermore/internal/dto"
	model "nevermore/internal/model/author"
	"nevermore/internal/storage"
)

type Service interface {
	Get(ctx context.Context, authorId int) (*dto.AuthorGetResponse, error)
	GetAuthorsList(ctx context.Context) ([]*dto.AuthorGetResponse, error)
	Update(ctx context.Context, author model.Author) error
	Delete(ctx context.Context, authorId int) error
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

func (s *service) Get(ctx context.Context, authorId int) (*dto.AuthorGetResponse, error) {
	author, err := s.st.DB().Author().Get(ctx, authorId)
	if err != nil {
		return author, fmt.Errorf("AuthorService:Get err -> %s", err.Error())
	}

	return author, nil
}

func (s *service) GetAuthorsList(ctx context.Context) ([]*dto.AuthorGetResponse, error) {
	result, err := s.st.DB().Author().GetAuthorsList(ctx)
	if err != nil {
		return nil, fmt.Errorf("AuthorService:GetAuthorsList err -> %s", err.Error())
	}

	return result, nil
}

func (s *service) Update(ctx context.Context, author model.Author) error {
	var err error

	err = s.st.DB().Author().Update(ctx, author)
	if err != nil {
		return fmt.Errorf("AuthorService:Update err -> %s", err.Error())
	}

	return nil
}

func (s *service) Delete(ctx context.Context, authorId int) error {
	err := s.st.DB().Author().Delete(ctx, authorId)
	if err != nil {
		return fmt.Errorf("AuthorService:Delete err -> %s", err.Error())
	}

	return nil
}
