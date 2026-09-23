package core_igdb_provider

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/M1sterZag/Dont_Play_Separately/internal/core/domain"
)

const platformFields string = "id,name,slug,abbreviation,platform_logo.url,checksum,updated_at"

func (c *Client) FetchPlatformsByIDs(ctx context.Context, ids []int64) ([]domain.Platform, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	idStr := strings.Join(IDsToString(ids), ",")

	query := fmt.Sprintf("fields %s; where id = (%s);", platformFields, idStr)

	raw, err := c.do(ctx, "platforms", query)
	if err != nil {
		return nil, err
	}

	var platforms []Platform
	if err := json.Unmarshal(raw, &platforms); err != nil {
		return nil, fmt.Errorf("unmarshal platforms: %w", err)
	}

	domainPlatforms := make([]domain.Platform, 0, len(platforms))
	for _, p := range platforms {
		domainPlatforms = append(domainPlatforms, p.ToDomain())
	}

	return domainPlatforms, nil
}

func (c *Client) SearchPlatforms(ctx context.Context, query string, limit, offset *int) ([]domain.Platform, error) {
	apicalypse := buildSearchQuery(query, platformFields, "", limit, offset)

	raw, err := c.do(ctx, "platforms", apicalypse)
	if err != nil {
		return nil, err
	}

	var platforms []Platform
	if err := json.Unmarshal(raw, &platforms); err != nil {
		return nil, fmt.Errorf("unmarshal search platforms: %w", err)
	}

	domainPlatforms := make([]domain.Platform, 0, len(platforms))
	for _, p := range platforms {
		domainPlatforms = append(domainPlatforms, p.ToDomain())
	}

	return domainPlatforms, nil
}
