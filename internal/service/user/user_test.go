package user_test

import (
	"context"
	//"errors"
	"nevermore/internal/dto"
	model "nevermore/internal/model/user"

	userService "nevermore/internal/service/user"
	minio2 "nevermore/internal/storage/minio"
	"nevermore/internal/storage/postgres"
	"nevermore/internal/storage/postgres/author"
	"nevermore/internal/storage/postgres/book"
	"nevermore/internal/storage/postgres/saved_author"
	"nevermore/internal/storage/postgres/user"
	"nevermore/internal/storage/redis"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockStorage остается без изменений
type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) DB() postgres.Repo {
	args := m.Called()
	return args.Get(0).(postgres.Repo)
}

func (m *MockStorage) Cash() redis.Storage {
	args := m.Called()
	return args.Get(0).(redis.Storage)
}

func (m *MockStorage) Cloud() minio2.Storage {
	args := m.Called()
	return args.Get(0).(minio2.Storage)
}

// MockRepo - полная реализация postgres.Repo
type MockRepo struct {
	mock.Mock
	userRepo        *MockUserStorage
	authorRepo      *MockAuthorStorage
	savedAuthorRepo *MockSavedAuthorStorage
	bookRepo        *MockBookStorage
}

func NewMockRepo(userRepo *MockUserStorage) *MockRepo {
	return &MockRepo{
		userRepo:        userRepo,
		authorRepo:      &MockAuthorStorage{},
		savedAuthorRepo: &MockSavedAuthorStorage{},
		bookRepo:        &MockBookStorage{},
	}
}

func (m *MockRepo) BeginTx(ctx context.Context) (*sqlx.Tx, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sqlx.Tx), args.Error(1)
}

func (m *MockRepo) User() user.Repo {
	return m.userRepo
}

func (m *MockRepo) Author() author.Repo {
	return nil
}

func (m *MockRepo) SavedAuthor() saved_author.Repo {
	return nil
}

func (m *MockRepo) Book() book.Repo {
	return nil
}

// MockUserStorage - только для пользователя
type MockUserStorage struct {
	mock.Mock
}

func (m *MockUserStorage) Get(ctx context.Context, id int) (*dto.UserGetResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.UserGetResponse), args.Error(1)
}

func (m *MockUserStorage) Update(ctx context.Context, user model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserStorage) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserStorage) Create(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserStorage) GetByEmail(ctx context.Context, email string) (model.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(model.User), args.Error(1)
}

// Заглушки для других репозиториев (они не используются в тестах пользователя)
type MockAuthorStorage struct {
	mock.Mock
}

func (m *MockAuthorStorage) Create(ctx context.Context) error {
	return nil
}

func (m *MockAuthorStorage) Get(ctx context.Context, id int) (*dto.AuthorGetResponse, error) {
	return nil, nil
}

func (m *MockAuthorStorage) Update(ctx context.Context) error {
	return nil
}

func (m *MockAuthorStorage) Delete(ctx context.Context, id int) error {
	return nil
}

type MockSavedAuthorStorage struct {
	mock.Mock
}

func (m *MockSavedAuthorStorage) Create(ctx context.Context) error {
	return nil
}

func (m *MockSavedAuthorStorage) Delete(ctx context.Context, id int) error {
	return nil
}

type MockBookStorage struct {
	mock.Mock
}

func (m *MockBookStorage) Create(ctx context.Context) error {
	return nil
}

func (m *MockBookStorage) Get(ctx context.Context, id int) error {
	return nil
}

func (m *MockBookStorage) Update(ctx context.Context) error {
	return nil
}

func (m *MockBookStorage) Delete(ctx context.Context, id int) error {
	return nil
}

// Mock CloudStorage
type MockCloudStorage struct {
	mock.Mock
}

func (m *MockCloudStorage) UploadPhoto(ctx context.Context, photo dto.FileInfo) (string, error) {
	args := m.Called(ctx, photo)
	return args.String(0), args.Error(1)
}

func (m *MockCloudStorage) DeletePhoto(ctx context.Context, photoURL string) error {
	args := m.Called(ctx, photoURL)
	return args.Error(0)
}

// setupTest создает моки для каждого теста
func setupTest(t *testing.T) (*MockStorage, *MockUserStorage, *MockCloudStorage, *MockRepo) {
	t.Helper()
	mockUserStorage := &MockUserStorage{}
	mockRepo := NewMockRepo(mockUserStorage)
	return &MockStorage{}, mockUserStorage, &MockCloudStorage{}, mockRepo
}

func TestUserService_Get(t *testing.T) {
	mockStorage, mockUserStorage, _, mockRepo := setupTest(t)

	ctx := context.Background()
	userID := 1
	expectedUser := &dto.UserGetResponse{
		Name:        "Test User",
		Email:       "test@example.com",
		PhoneNumber: "+1234567890",
		Photo:       nil,
	}

	mockStorage.On("DB").Return(mockRepo)
	mockUserStorage.On("Get", ctx, userID).Return(expectedUser, nil)

	service := userService.New(mockStorage)
	result, err := service.Get(ctx, userID)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser.Name, result.Name)
	assert.Equal(t, expectedUser.Email, result.Email)
	assert.Equal(t, expectedUser.PhoneNumber, result.PhoneNumber)

	mockStorage.AssertExpectations(t)
	mockUserStorage.AssertExpectations(t)
}
