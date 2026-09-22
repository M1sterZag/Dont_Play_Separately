package core_http_request

import (
	"fmt"
	"net/http"
	"strconv"

	core_errors "github.com/M1sterZag/Dont_Play_Separately/internal/core/errors"
	"github.com/google/uuid"
)

func GetUUIDPathParam(r *http.Request, key string) (uuid.UUID, error) {
	pathValue := r.PathValue(key)
	if pathValue == "" {
		return uuid.UUID{}, fmt.Errorf("no key='%s' in path values: %w", key, core_errors.ErrInvalidArgument)
	}

	val, err := uuid.Parse(pathValue)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("path value='%s' by key='%s' not a valid uuid: %w", pathValue, key, core_errors.ErrInvalidArgument)
	}

	return val, nil
}

func GetInt64PathParam(r *http.Request, key string) (int64, error) {
	pathValue := r.PathValue(key)
	if pathValue == "" {
		return 0, fmt.Errorf("no key='%s' in path values: %w", key, core_errors.ErrInvalidArgument)
	}

	val, err := strconv.ParseInt(pathValue, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("path value='%s' by key='%s' not a valid int64: %w", pathValue, key, core_errors.ErrInvalidArgument)
	}
	return val, nil
}
