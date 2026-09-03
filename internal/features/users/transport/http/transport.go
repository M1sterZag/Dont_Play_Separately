package users_transport_http

import (
	"context"
	"net/http"

	"github.com/M1sterZag/Dont_Play_Separately/internal/core/domain"
	core_storage "github.com/M1sterZag/Dont_Play_Separately/internal/core/storage"
	core_http_server "github.com/M1sterZag/Dont_Play_Separately/internal/core/transport/server"
	"github.com/google/uuid"
)

type UsersHTTPHandler struct {
	usersService UsersService
	storage      core_storage.Storage
}

type UsersService interface {
	GetProfileByID(ctx context.Context, userID uuid.UUID) (domain.UserProfile, error)
}

func NewUsersHTTPHandler(userService UsersService, storage core_storage.Storage) *UsersHTTPHandler {
	return &UsersHTTPHandler{
		usersService: userService,
		storage:      storage,
	}
}

func (h *UsersHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodGet,
			Path:    "/users/profile/{user_id}",
			Handler: h.GetProfileByID,
		},
	}
}
