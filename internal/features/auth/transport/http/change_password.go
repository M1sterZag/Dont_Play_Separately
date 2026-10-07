package auth_transport_http

import (
	"net/http"

	core_errors "github.com/M1sterZag/Dont_Play_Separately/internal/core/errors"
	core_logger "github.com/M1sterZag/Dont_Play_Separately/internal/core/logger"
	core_http_middleware "github.com/M1sterZag/Dont_Play_Separately/internal/core/transport/http/middleware"
	core_http_request "github.com/M1sterZag/Dont_Play_Separately/internal/core/transport/http/request"
	core_http_response "github.com/M1sterZag/Dont_Play_Separately/internal/core/transport/http/response"
)

// ChangePassword changes the authenticated user's password.
// @Summary Change password
// @Description Updates the authenticated user's password after verifying the current one. All active refresh sessions are revoked, so the user must log in again on all devices.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body ChangePasswordRequest true "Old and new passwords"
// @Success 204 "No Content"
// @Failure 400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure 401 {object} core_http_response.ErrorResponse "Unauthorized"
// @Failure 403 {object} core_http_response.ErrorResponse "Forbidden (wrong current password)"
// @Failure 500 {object} core_http_response.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /auth/change-password [patch]
func (h *AuthHTTPHandler) ChangePassword(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	userID, ok := core_http_middleware.UserIDFromContext(ctx)
	if !ok {
		responseHandler.ErrorResponse(core_errors.ErrUnauthenticated, "failed to get user_id (context)")
		return
	}

	var request ChangePasswordRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate request body")
		return
	}

	if err := h.authService.ChangePassword(ctx, userID, request.OldPassword, request.NewPassword); err != nil {
		responseHandler.ErrorResponse(err, "failed to change password")
		return
	}

	responseHandler.NoContentResponse()
}