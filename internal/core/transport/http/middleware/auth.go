package core_http_middleware

import (
	"context"
	"net/http"
	"strings"

	core_errors "github.com/M1sterZag/Dont_Play_Separately/internal/core/errors"
	core_logger "github.com/M1sterZag/Dont_Play_Separately/internal/core/logger"
	core_http_response "github.com/M1sterZag/Dont_Play_Separately/internal/core/transport/http/response"
	"github.com/google/uuid"
)

type contextKey string

const (
	userUUIDContextKey contextKey = "user_uuid"
	bearerPrefix       string     = "Bearer "
	authorizationKey   string     = "Authorization"
)

func Auth(parseToken func(token string) (uuid.UUID, error)) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			responseHandler := core_http_response.NewHTTPResponseHandler(
				core_logger.FromContext(r.Context()), w,
			)

			header := r.Header.Get(authorizationKey)
			if !strings.HasPrefix(header, bearerPrefix) {
				responseHandler.ErrorResponse(core_errors.ErrUnauthenticated, "missing bearer token")
				return
			}

			token := strings.TrimPrefix(header, bearerPrefix)
			userUUID, err := parseToken(token)
			if err != nil {
				responseHandler.ErrorResponse(err, "invalid access token")
				return
			}

			ctx := context.WithValue(r.Context(), userUUIDContextKey, userUUID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func UserUUIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	val, ok := ctx.Value(userUUIDContextKey).(uuid.UUID)
	return val, ok
}
