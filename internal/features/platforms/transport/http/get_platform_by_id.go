package platforms_transport_http

import (
	"net/http"

	core_logger "github.com/M1sterZag/Dont_Play_Separately/internal/core/logger"
	core_http_request "github.com/M1sterZag/Dont_Play_Separately/internal/core/transport/http/request"
	core_http_response "github.com/M1sterZag/Dont_Play_Separately/internal/core/transport/http/response"
)

type GetPlatformResponse PlatformDTOResponse

// GetPlatformByID returns a platform by its IGDB ID.
// @Summary Get platform by ID
// @Description Returns a catalog platform by its numeric ID.
// @Tags platforms
// @Produce json
// @Param id path int true "Platform ID"
// @Success 200 {object} GetPlatformResponse "OK"
// @Failure 400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure 404 {object} core_http_response.ErrorResponse "Not found"
// @Failure 500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router /platforms/{id} [get]
func (h *PlatformsHTTPHandler) GetPlatformByID(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	platformID, err := core_http_request.GetInt64PathParam(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get platform id path param")
		return
	}

	platform, err := h.platformsService.GetPlatformByID(ctx, platformID)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get platform")
		return
	}

	response := GetPlatformResponse(platformDTOFromDomain(platform))
	responseHandler.JSONResponse(response, http.StatusOK)
}