package auth_transport_http

import (
	"net/http"

	core_logger "github.com/M1sterZag/Dont_Play_Separately/internal/core/logger"
	core_http_request "github.com/M1sterZag/Dont_Play_Separately/internal/core/transport/http/request"
	core_http_response "github.com/M1sterZag/Dont_Play_Separately/internal/core/transport/http/response"
)

// Login authenticates a user.
// @Summary Login
// @Description Authenticates a user by email and password and returns an access and refresh token pair.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "User credentials"
// @Success 200 {object} TokenResponse "OK"
// @Failure 400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure 401 {object} core_http_response.ErrorResponse "Unauthorized"
// @Failure 500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router /auth/login [post]
func (h *AuthHTTPHandler) Login(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	var request LoginRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")
		return
	}

	tokens, err := h.authService.Login(ctx, request.Email, request.Password)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to login user")
		return
	}

	responseHandler.JSONResponse(TokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, http.StatusOK)
}
