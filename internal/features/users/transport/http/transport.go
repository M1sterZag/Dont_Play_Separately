package users_transport_http

import (
	"context"
	"net/http"

	"github.com/M1sterZag/Dont_Play_Separately/internal/core/domain"
	core_http_server "github.com/M1sterZag/Dont_Play_Separately/internal/core/transport/server"
	"github.com/google/uuid"
)

type UsersHTTPHandler struct {
	usersService UsersService
}

type UsersService interface {
	GetUserByUUID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	PatchUser(ctx context.Context, profilePatch domain.User) (domain.User, error)
}

func NewUsersHTTPHandler(userService UsersService) *UsersHTTPHandler {
	return &UsersHTTPHandler{
		usersService: userService,
	}
}

func (h *UsersHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodGet,
			Path:    "/users/profile",
			Handler: h.GetUserByUUID,
		},
		{
			Method:  http.MethodPatch,
			Path:    "/users/profile",
			Handler: h.PatchUser,
		},
	}
}
