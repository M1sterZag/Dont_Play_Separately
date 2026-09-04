package users_transport_http

import (
	"net/http"

	core_errors "github.com/M1sterZag/Dont_Play_Separately/internal/core/errors"
	core_logger "github.com/M1sterZag/Dont_Play_Separately/internal/core/logger"
	core_http_middleware "github.com/M1sterZag/Dont_Play_Separately/internal/core/transport/http/middleware"
	core_http_response "github.com/M1sterZag/Dont_Play_Separately/internal/core/transport/http/response"
)

// DeleteProfile deletes the authenticated user's account.
// @Summary Delete user profile
// @Description Deletes the authenticated user's profile (account). Only the owner can delete their own profile.
// @Tags users
// @Produce json
// @Success 204 "No Content"
// @Failure 401 {object} core_http_response.ErrorResponse "Unauthorized"
// @Failure 404 {object} core_http_response.ErrorResponse "Not found"
// @Failure 500 {object} core_http_response.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /users/profile/me [delete]
func (h *UsersHTTPHandler) DeleteProfile(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	userID, ok := core_http_middleware.UserIDFromContext(ctx)
	if !ok {
		responseHandler.ErrorResponse(core_errors.ErrUnauthenticated, "failed to get user id from context")
		return
	}
	if err := h.usersService.DeleteProfile(ctx, userID); err != nil {
		responseHandler.ErrorResponse(err, "failed to delete profile")
		return
	}

	responseHandler.NoContentResponse()
}
