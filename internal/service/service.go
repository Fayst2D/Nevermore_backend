package service

import (
	"nevermore/internal/service/author"
	"nevermore/internal/service/authorization"
	"nevermore/internal/service/book"
	"nevermore/internal/service/chat"
	"nevermore/internal/service/saved_author"
	"nevermore/internal/service/user"
	"nevermore/internal/storage"
	"nevermore/pkg/auth"

	"github.com/gammazero/workerpool"

	"nevermore/pkg/hash"
)

type Service interface {
	User() user.Service
	Author() author.Service
	SavedAuthor() saved_author.Service
	Book() book.Service
	Auth() authorization.Service
	Chat() chat.Service
}

type service struct {
	user        user.Service
	author      author.Service
	savedAuthor saved_author.Service
	book        book.Service
	auth        authorization.Service
	chat        chat.Service
}

func New(st storage.Storage,
	hash hash.PasswordHasher,
	manager auth.TokenManager,
	wp *workerpool.WorkerPool) Service {

	//go chat.Run()
	result := &service{
		user:   user.New(st),
		author: author.New(st),
		book:   book.New(st, wp),
		auth:   authorization.New(st, manager, hash),
		chat:   chat.New(st, wp),
	}

	return result
}

func (s *service) User() user.Service {
	return s.user
}
func (s *service) Author() author.Service {
	return s.author
}
func (s *service) SavedAuthor() saved_author.Service {
	return s.savedAuthor
}
func (s *service) Book() book.Service {
	return s.book
}
func (s *service) Auth() authorization.Service {
	return s.auth
}

func (s *service) Chat() chat.Service {
	return s.chat
}
