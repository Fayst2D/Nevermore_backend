package authorization

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"nevermore/pkg/logger"
	"strconv"
	"time"

	"nevermore/internal/dto"
	"nevermore/internal/storage"
	"nevermore/pkg/auth"
	"nevermore/pkg/hash"

	model "nevermore/internal/model/user"

	"github.com/asaskevich/govalidator"
)

type Service interface {
	Register(ctx context.Context, user *dto.RegisterRequest) error
	Login(ctx context.Context, user *dto.LoginRequest) (*auth.Token, error)
	Logout(ctx context.Context, email string) error
	Refresh(ctx context.Context, rt string) (*auth.Token, error)
}

type service struct {
	st storage.Storage

	manager auth.TokenManager
	hash    hash.PasswordHasher
}

func New(st storage.Storage, manager auth.TokenManager, hash hash.PasswordHasher) Service {

	result := &service{
		st:      st,
		manager: manager,
		hash:    hash,
	}

	return result
}

func (s *service) Register(ctx context.Context, req *dto.RegisterRequest) error {

	if !govalidator.IsEmail(req.Email) {
		return fmt.Errorf("invalid email")
	}

	userInfo, err := s.st.DB().User().GetByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("AuthService:Register err -> %s", err.Error())
	}

	if userInfo.Id > 0 {
		return errors.New("user already exists")
	}

	req.Password, err = s.hash.Hash(req.Password)
	if err != nil {
		return err
	}

	user := &model.User{
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
		Email:       req.Email,
		Role:        "user",
		Password:    req.Password,
		Photo:       nil,
		CreatedAt:   time.Now().UTC(),
	}

	err = s.st.DB().User().Create(ctx, user)
	if err != nil {
		return fmt.Errorf("AuthService:Register err -> %s", err.Error())
	}
	return nil
}

func (s *service) Login(ctx context.Context, req *dto.LoginRequest) (*auth.Token, error) {
	var err error
	userInfo, err := s.st.DB().User().GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("AuthService:Login err -> %s", err.Error())
	}

	passwordHash, err := s.hash.Hash(req.Password)
	if err != nil {
		return nil, err
	}

	if passwordHash != userInfo.Password {
		return nil, errors.New("invalid password")
	}

	idStr := strconv.Itoa(userInfo.Id)

	token, err := s.manager.GenerateTokenPair(idStr)
	if err != nil {
		return nil, fmt.Errorf("AuthService:Login err -> %s", err.Error())
	}

	log := logger.Get()
	log.Info().Msgf("login success -> %s, %d", token, userInfo.Id)

	err = s.st.Cash().Set(ctx, fmt.Sprintf("refresh_token:%s", userInfo.Email), token.Rt, 24*30*time.Hour)

	return token, err
}

func (s *service) Logout(ctx context.Context, email string) error {
	err := s.st.Cash().Delete(ctx, fmt.Sprintf("refresh_token:%s", email))
	if err != nil {
		return fmt.Errorf("AuthService:Logout err -> %s", err.Error())
	}

	return nil
}

func (s *service) Refresh(ctx context.Context, rt string) (*auth.Token, error) {
	//idStr, err := s.st.Cash().Get(ctx, rt)
	//if err != nil {
	//	return nil, fmt.Errorf("AuthService:Refresh err -> %s", err.Error())
	//}
	//
	//err = s.st.Cash().Delete(ctx, rt)
	//if err != nil {
	//	return nil, fmt.Errorf("AuthService:Refresh err -> %s", err.Error())
	//}
	//
	//id, err := strconv.Atoi(idStr)
	//if err != nil {
	//	return nil, fmt.Errorf("AuthService:Refresh err -> %s", err.Error())
	//}

	tokens, err := s.manager.RefreshAccessToken(rt)
	if err != nil {
		return nil, fmt.Errorf("AuthService:Refresh -> %s", err.Error())
	}

	return tokens, nil
}
