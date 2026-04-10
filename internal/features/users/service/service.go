package users_service

import (
	"context"

	"github.com/annpolukaro/to_do_app/internal/core/domain"
)

type UsersService struct {
	usersRepository usersRepository
}

type usersRepository interface {
	CreateUser(
		ctx context.Context,
		user domain.User,

	) (domain.User, error)
}

func NewUsersService(
	usersRepository usersRepository,
) *UsersService {
	return &UsersService{
		usersRepository: usersRepository,
	}
}
