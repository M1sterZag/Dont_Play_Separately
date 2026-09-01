package auth_transport_http

import (
	"context"
	"net/http"

	core_http_server "github.com/M1sterZag/Dont_Play_Separately/internal/core/transport/server"
	auth_service "github.com/M1sterZag/Dont_Play_Separately/internal/features/auth/service"
)

type AuthService interface {
	Register(ctx context.Context, email, password, nickname string) (auth_service.Tokens, error)
	Login(ctx context.Context, email, password string) (auth_service.Tokens, error)
	Refresh(ctx context.Context, refreshToken string) (auth_service.Tokens, error)
	Logout(ctx context.Context, refreshToken string) error
}

type AuthHTTPHandler struct {
	authService AuthService
}

func NewAuthHTTPHandler(authService AuthService) *AuthHTTPHandler {
	return &AuthHTTPHandler{
		authService: authService,
	}
}

func (h *AuthHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/auth/register",
			Handler: h.Register,
		},
		{
			Method:  http.MethodPost,
			Path:    "/auth/login",
			Handler: h.Login,
		},
		{
			Method:  http.MethodPost,
			Path:    "/auth/logout",
			Handler: h.Logout,
		},
		{
			Method:  http.MethodPost,
			Path:    "/auth/refresh",
			Handler: h.Refresh,
		},
	}
}
