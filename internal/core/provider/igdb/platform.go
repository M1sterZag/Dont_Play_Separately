package core_igdb_provider

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

const platformFields string = "id,name,slug,abbreviation,platform_logo.url,checksum,updated_at"

func (c *Client) FetchPlatformsByIDs(ctx context.Context, ids []int64) ([]Platform, error) {
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

	return platforms, nil
}

func (c *Client) SearchPlatforms(ctx context.Context, query string, limit int) ([]Platform, error) {
	apicalypse := fmt.Sprintf(
		"search \"%s\"; fields %s; limit %d;",
		query, platformFields, limit,
	)

	raw, err := c.do(ctx, "platforms", apicalypse)
	if err != nil {
		return nil, err
	}

	var platforms []Platform
	if err := json.Unmarshal(raw, &platforms); err != nil {
		return nil, fmt.Errorf("unmarshal search platforms: %w", err)
	}

	return platforms, nil
}
