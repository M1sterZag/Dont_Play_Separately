package users_transport_http

import (
	"net/http"

	core_logger "github.com/M1sterZag/Dont_Play_Separately/internal/core/logger"
	core_http_request "github.com/M1sterZag/Dont_Play_Separately/internal/core/transport/http/request"
	core_http_response "github.com/M1sterZag/Dont_Play_Separately/internal/core/transport/http/response"
)

type GetProfileResponse UserProfileDTOResponse

func (h *UsersHTTPHandler) GetProfileByID(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	userID, err := core_http_request.GetUUIDPathParam(r, "user_id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get user id path param")
		return
	}

	profile, err := h.usersService.GetProfileByID(ctx, userID)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get user profile")
		return
	}

	response := GetProfileResponse(userProfileDTOFromDomain(profile))
	responseHandler.JSONResponse(response, http.StatusOK)
}
