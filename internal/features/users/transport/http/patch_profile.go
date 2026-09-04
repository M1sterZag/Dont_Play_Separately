package users_transport_http

import (
	"net/http"

	core_errors "github.com/M1sterZag/Dont_Play_Separately/internal/core/errors"
	core_logger "github.com/M1sterZag/Dont_Play_Separately/internal/core/logger"
	core_http_middleware "github.com/M1sterZag/Dont_Play_Separately/internal/core/transport/http/middleware"
	core_http_request "github.com/M1sterZag/Dont_Play_Separately/internal/core/transport/http/request"
	core_http_response "github.com/M1sterZag/Dont_Play_Separately/internal/core/transport/http/response"
)

type PatchProfileResponse UserProfileDTOResponse

// PatchProfile updates the authenticated user's profile.
// @Summary Update user profile
// @Description Partially updates the authenticated user's public profile (nickname, bio, avatar preset). A field is updated only if it is present in the body; pass "bio": null to clear it. Only the owner can modify their own profile.
// @Tags users
// @Accept json
// @Produce json
// @Param request body PatchProfileRequest true "Profile patch payload"
// @Success 200 {object} PatchProfileResponse "OK"
// @Failure 400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure 401 {object} core_http_response.ErrorResponse "Unauthorized"
// @Failure 404 {object} core_http_response.ErrorResponse "Not found"
// @Failure 409 {object} core_http_response.ErrorResponse "Conflict (optimistic lock)"
// @Failure 500 {object} core_http_response.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /users/profile/me [patch]
func (h *UsersHTTPHandler) PatchProfile(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	userID, ok := core_http_middleware.UserIDFromContext(ctx)
	if !ok {
		responseHandler.ErrorResponse(core_errors.ErrUnauthenticated, "failed to get user id from context")
		return
	}

	var request PatchProfileRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")
		return
	}

	profilePatch := userProfilePatchFromRequest(request)
	profileDomain, err := h.usersService.PatchProfile(ctx, userID, profilePatch)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to patch profile")
		return
	}

	response := PatchProfileResponse(userProfileDTOFromDomain(profileDomain, h.storage))
	responseHandler.JSONResponse(response, http.StatusOK)

}
