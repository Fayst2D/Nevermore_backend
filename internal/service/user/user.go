package user

import (
	"context"
	"fmt"
	"nevermore/internal/dto"
	model "nevermore/internal/model/user"
	"nevermore/internal/storage"
)

type Service interface {
	Get(ctx context.Context, userId int) (*dto.UserGetResponse, error)
	Update(ctx context.Context, userId int, req dto.UpdateUserRequest, photo dto.FileInfo) error
	Delete(ctx context.Context, userId int) error
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

func (s *service) Get(ctx context.Context, userId int) (*dto.UserGetResponse, error) {
	user, err := s.st.DB().User().Get(ctx, userId)
	if err != nil {
		return user, fmt.Errorf("UserService:Get err -> %s", err.Error())
	}

	return user, nil
}

func (s *service) Update(ctx context.Context, userId int, req dto.UpdateUserRequest, photo dto.FileInfo) error {
	var err error

	user := model.User{
		Id:          userId,
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
		Email:       req.Email,
	}

	if photo.File != nil {
		photoUrl, err := s.st.Cloud().UploadPhoto(ctx, photo)
		user.Photo = &photoUrl
		if err != nil {
			return fmt.Errorf("UserService:Update err -> %s", err.Error())
		}
	}

	err = s.st.DB().User().Update(ctx, user)
	if err != nil {
		return fmt.Errorf("UserService:Update err -> %s", err.Error())
	}

	return nil
}

func (s *service) Delete(ctx context.Context, userId int) error {
	err := s.st.DB().User().Delete(ctx, userId)
	if err != nil {
		return fmt.Errorf("UserService:Delete err -> %s", err.Error())
	}

	return nil
}
