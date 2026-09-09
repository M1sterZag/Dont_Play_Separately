package core_igdb_provider

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

const gameFields string = "id,name,slug,cover.url,checksum,updated_at"

func (c *Client) FetchGamesByIDs(ctx context.Context, ids []int64) ([]Game, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	idStr := strings.Join(IDsToString(ids), ",")

	query := fmt.Sprintf("fields %s; where id = (%s);", gameFields, idStr)

	raw, err := c.do(ctx, "games", query)
	if err != nil {
		return nil, err
	}

	var games []Game
	if err := json.Unmarshal(raw, &games); err != nil {
		return nil, fmt.Errorf("unmarshal games: %w", err)
	}

	return games, nil
}

func (c *Client) SearchGames(ctx context.Context, query string, limit int) ([]Game, error) {
	apicalypse := fmt.Sprintf(
		"search \"%s\"; fields %s; limit %d;",
		query, gameFields, limit,
	)

	raw, err := c.do(ctx, "games", apicalypse)
	if err != nil {
		return nil, err
	}

	var games []Game
	if err := json.Unmarshal(raw, &games); err != nil {
		return nil, fmt.Errorf("unmarshal search games: %w", err)
	}

	return games, nil
}
