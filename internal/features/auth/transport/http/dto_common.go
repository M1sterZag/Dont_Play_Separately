package auth_transport_http

type RegisterRequest struct {
	Email               string `json:"email" validate:"required,email" example:"user@example.com"`
	Password            string `json:"password" validate:"required,password" example:"P@ssw0rd123!"`
	Nickname            string `json:"nickname" validate:"required,min=1,max=40" example:"gamemaster"`
	FavoritePlatformIDs []int  `json:"favorite_platform_ids" validate:"omitempty,max=3,dive,gt=0"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email" example:"user@example.com"`
	Password string `json:"password" validate:"required" example:"P@ssw0rd123!"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

type VerifyEmailRequest struct {
	Email string `json:"email" validate:"required,email" example:"user@example.com"`
	Code  string `json:"code" validate:"required,len=6,numeric" example:"123456"`
}

type ResendVerificationCodeRequest struct {
	Email string `json:"email" validate:"required,email" example:"user@example.com"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
