package users_transport_http

import (
	"net/http"

	core_logger "github.com/M1sterZag/Dont_Play_Separately/internal/core/logger"
	core_http_request "github.com/M1sterZag/Dont_Play_Separately/internal/core/transport/http/request"
	core_http_response "github.com/M1sterZag/Dont_Play_Separately/internal/core/transport/http/response"
)

type GetUserResponse UserDTOResponse

func (h *UsersHTTPHandler) GetUserByUUID(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	userUUID, err := core_http_request.GetUUIDPathParam(r, "user_uuid")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get user uuid path param")
		return
	}

	userDomain, err := h.usersService.GetUserByUUID(ctx, userUUID)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get user")
		return
	}

	response := GetUserResponse(userDTOFromDomain(userDomain))
	responseHandler.JSONResponse(response, http.StatusOK)
}
