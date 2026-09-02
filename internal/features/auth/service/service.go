package auth_service

import (
	"context"

	"github.com/M1sterZag/Dont_Play_Separately/internal/core/domain"
	"github.com/google/uuid"
)

type AuthRepository interface {
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
	CreateUser(ctx context.Context, user domain.User) (domain.User, error)

	CreateSession(ctx context.Context, session domain.RefreshSession) error
	FindSessionByUUID(ctx context.Context, sessionUUID uuid.UUID) (domain.RefreshSession, error)
	RevokeSession(ctx context.Context, sessionUUID uuid.UUID) error
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type AuthService struct {
	authRepository AuthRepository
	jwtSigner      *JWTSigner
}

func NewAuthService(authRepository AuthRepository, jwtSigner *JWTSigner) *AuthService {
	return &AuthService{
		authRepository: authRepository,
		jwtSigner:      jwtSigner,
	}
}
