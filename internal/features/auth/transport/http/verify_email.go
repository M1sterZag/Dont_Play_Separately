package auth_transport_http

import (
	"net/http"

	core_logger "github.com/M1sterZag/Dont_Play_Separately/internal/core/logger"
	core_http_request "github.com/M1sterZag/Dont_Play_Separately/internal/core/transport/http/request"
	core_http_response "github.com/M1sterZag/Dont_Play_Separately/internal/core/transport/http/response"
)

// VerifyEmail verifies the user email with the code sent during registration.
// @Summary Verify user email with a code
// @Description Confirms the email using the verification code sent during registration. On success the account is marked as verified and an access and refresh token pair is returned.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body VerifyEmailRequest true "Verification payload"
// @Success 200 {object} TokenResponse "Verified, tokens returned"
// @Failure 400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure 404 {object} core_http_response.ErrorResponse "Not found"
// @Failure 409 {object} core_http_response.ErrorResponse "Already verified"
// @Failure 500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router /auth/verify-email [post]
func (h *AuthHTTPHandler) VerifyEmail(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	var request VerifyEmailRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate request body")
		return
	}

	tokens, err := h.authService.VerifyEmail(ctx, request.Email, request.Code)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to verify email")
		return
	}

	responseHandler.JSONResponse(TokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, http.StatusOK)
}