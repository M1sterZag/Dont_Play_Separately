package core_http_request

import (
	"fmt"
	"net/http"

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
