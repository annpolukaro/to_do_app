package user_transport_http

import (
	"context"
	"net/http"

	"github.com/annpolukaro/to_do_app/internal/core/domain"
	core_http_server "github.com/annpolukaro/to_do_app/internal/core/transport/http/server"
)

type UserHTTPHandler struct {
	usersService UsersService
}

type UsersService interface {
	CreateUser(
		ctx context.Context,
		user domain.User,

	) (domain.User, error)
}

func NewUsersHTTPHandler(UsersService UsersService) *UserHTTPHandler {
	return &UserHTTPHandler{
		usersService: UsersService,
	}
}

func (h *UserHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/users",
			Handler: h.CreateUser,
		},
	}
}
