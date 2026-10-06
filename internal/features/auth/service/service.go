package auth_service

import (
	"context"

	core_cache "github.com/M1sterZag/Dont_Play_Separately/internal/core/cache"
	"github.com/M1sterZag/Dont_Play_Separately/internal/core/domain"
	core_mailer "github.com/M1sterZag/Dont_Play_Separately/internal/core/mailer"
	auth_config "github.com/M1sterZag/Dont_Play_Separately/internal/features/auth"
	"github.com/google/uuid"
)

type AuthRepository interface {
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
	CreateUser(ctx context.Context, user domain.User, favoritePlatformIDs []int) (domain.User, error)
	MarkEmailVerified(ctx context.Context, userID uuid.UUID) error

	CreateSession(ctx context.Context, session domain.RefreshSession) error
	FindSessionByID(ctx context.Context, sessionID uuid.UUID) (domain.RefreshSession, error)
	RevokeSession(ctx context.Context, sessionID uuid.UUID) error
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type AuthService struct {
	authRepository AuthRepository
	jwtSigner      *JWTSigner
	cache          core_cache.Cache
	mailer         core_mailer.Mailer
	config         auth_config.Config
}

func NewAuthService(
	authRepository AuthRepository,
	jwtSigner *JWTSigner,
	cache core_cache.Cache,
	mailer core_mailer.Mailer,
	config auth_config.Config,
) *AuthService {
	return &AuthService{
		authRepository: authRepository,
		jwtSigner:      jwtSigner,
		cache:          cache,
		mailer:         mailer,
		config:         config,
	}
}
