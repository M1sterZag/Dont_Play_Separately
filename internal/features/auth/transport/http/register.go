package auth_transport_http

import (
	"net/http"

	core_logger "github.com/M1sterZag/Dont_Play_Separately/internal/core/logger"
	core_http_request "github.com/M1sterZag/Dont_Play_Separately/internal/core/transport/http/request"
	core_http_response "github.com/M1sterZag/Dont_Play_Separately/internal/core/transport/http/response"
)

// Register creates a new user.
// @Summary Register a new user
// @Description Registers a new user by email, password and nickname and returns an access and refresh token pair.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Registration payload"
// @Success 201 {object} TokenResponse "Created"
// @Failure 400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure 409 {object} core_http_response.ErrorResponse "Conflict"
// @Failure 500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router /auth/register [post]
func (h *AuthHTTPHandler) Register(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	var request RegisterRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")
		return
	}

	tokens, err := h.authService.Register(ctx, request.Email, request.Password, request.Nickname)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to register")
		return
	}

	responseHandler.JSONResponse(TokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, http.StatusCreated)
}
