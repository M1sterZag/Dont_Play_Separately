package auth_transport_http

import (
	"net/http"

	core_logger "github.com/M1sterZag/Dont_Play_Separately/internal/core/logger"
	core_http_request "github.com/M1sterZag/Dont_Play_Separately/internal/core/transport/http/request"
	core_http_response "github.com/M1sterZag/Dont_Play_Separately/internal/core/transport/http/response"
)

// ResendVerificationCode resends the verification code to the user email.
// @Summary Resend verification code
// @Description Generates a new verification code and sends it to the user email. Can be requested no more often than once per configured interval.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body ResendVerificationCodeRequest true "Resend payload"
// @Success 202 {object} MessageResponse "Accepted, new code sent"
// @Failure 400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure 404 {object} core_http_response.ErrorResponse "User not found"
// @Failure 409 {object} core_http_response.ErrorResponse "Already verified or too frequent requests"
// @Failure 500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router /auth/resend-verification-code [post]
func (h *AuthHTTPHandler) ResendVerificationCode(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	var request ResendVerificationCodeRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate request body")
		return
	}

	if err := h.authService.ResendVerificationCode(ctx, request.Email); err != nil {
		responseHandler.ErrorResponse(err, "failed to resend verification code")
		return
	}

	responseHandler.JSONResponse(MessageResponse{
		Message: "verification code sent to email",
	}, http.StatusAccepted)
}