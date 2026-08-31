package auth_transport_http

import (
	"context"

	auth_service "github.com/M1sterZag/Dont_Play_Separately/internal/features/auth/service"
)

type AuthService interface {
	Register(ctx context.Context, email, password, nickname string) (auth_service.Tokens, error)
	Login(ctx context.Context, email, password string) (auth_service.Tokens, error)
	Refresh(ctx context.Context, refreshToken string) (auth_service.Tokens, error)
	Logout(ctx context.Context, refreshToken string) error
}
