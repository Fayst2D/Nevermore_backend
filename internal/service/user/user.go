package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"nevermore/internal/dto"
	model "nevermore/internal/model/user"
	"nevermore/internal/storage"
	"nevermore/pkg/email"

	"github.com/google/uuid"
)

type Service interface {
	Get(ctx context.Context, userId int) (*dto.UserGetResponse, error)
	Update(ctx context.Context, userId int, req dto.UpdateUserRequest, photo dto.FileInfo) error
	Delete(ctx context.Context, userId int) error
}

type service struct {
	st       storage.Storage
	emailSrv email.Service
}

func New(st storage.Storage, emailSrv email.Service) Service {
	result := &service{
		st:       st,
		emailSrv: emailSrv,
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

	// Проверяем, не используется ли email другим пользователем
	currentUser, err := s.st.DB().User().GetByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("UserService:Update err -> %s", err.Error())
	}
	if err == nil && currentUser.Id != userId {
		// Email уже используется другим пользователем
		return fmt.Errorf("email already in use")
	}

	// Получаем пользователя по ID для проверки изменения email
	existingUser, err := s.st.DB().User().GetById(ctx, userId)
	if err != nil {
		return fmt.Errorf("UserService:Update err -> %s", err.Error())
	}

	user := model.User{
		Id:          userId,
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
		Email:       req.Email,
	}

	// Если email изменился, генерируем новый токен подтверждения
	if existingUser.Email != req.Email {
		verificationToken := uuid.New().String()
		err = s.st.DB().User().UpdateVerificationToken(ctx, req.Email, verificationToken)
		if err != nil {
			return fmt.Errorf("UserService:Update err -> %s", err.Error())
		}

		// Отправляем email с новым токеном подтверждения
		err = s.emailSrv.SendVerificationEmail(req.Email, verificationToken)
		if err != nil {
			// Логируем ошибку, но не прерываем обновление
			// Пользователь сможет запросить повторную отправку позже
		}
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
