package auth_service

import auth_repository "github.com/M1sterZag/Dont_Play_Separately/internal/features/auth/repository"

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type AuthService struct {
	authRepository auth_repository.AuthRepository
	jwtSigner      *JWTSigner
}

func NewAuthService(authRepository auth_repository.AuthRepository, jwtSigner *JWTSigner) *AuthService {
	return &AuthService{
		authRepository: authRepository,
		jwtSigner:      jwtSigner,
	}
}
