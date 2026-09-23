package platforms_transport_http

import (
	"context"
	"net/http"

	"github.com/M1sterZag/Dont_Play_Separately/internal/core/domain"
	core_http_server "github.com/M1sterZag/Dont_Play_Separately/internal/core/transport/server"
)

type PlatformsService interface {
	GetPlatformByID(ctx context.Context, platformID int64) (domain.Platform, error)
	SearchPlatforms(ctx context.Context, searchQuery string, limit, offset *int) ([]domain.Platform, error)
}

type PlatformsHTTPHandler struct {
	platformsService PlatformsService
}

func NewPlatformsHTTPHandler(platformsService PlatformsService) *PlatformsHTTPHandler {
	return &PlatformsHTTPHandler{
		platformsService: platformsService,
	}
}

func (h *PlatformsHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodGet,
			Path:    "/platforms/{id}",
			Handler: h.GetPlatformByID,
		},
		{
			Method:  http.MethodGet,
			Path:    "/platforms/search",
			Handler: h.SearchPlatforms,
		},
	}
}
