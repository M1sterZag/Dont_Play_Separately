package core_igdb_provider

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/M1sterZag/Dont_Play_Separately/internal/core/domain"
)

const gameFields string = "id,name,slug,cover.url,game_type,checksum,updated_at"

const gameTypeMainGame = 0

func (c *Client) FetchGamesByIDs(ctx context.Context, ids []int64) ([]domain.Game, error) {
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

	domainGames := make([]domain.Game, 0, len(games))
	for _, g := range games {
		domainGames = append(domainGames, g.ToDomain())
	}

	return domainGames, nil
}

func (c *Client) SearchGames(ctx context.Context, query string, limit, offset *int) ([]domain.Game, error) {
	apicalypse := buildSearchQuery(query, gameFields, "version_parent = null & game_type = 0", limit, offset)

	raw, err := c.do(ctx, "games", apicalypse)
	if err != nil {
		return nil, err
	}

	var games []Game
	if err := json.Unmarshal(raw, &games); err != nil {
		return nil, fmt.Errorf("unmarshal search games: %w", err)
	}

	domainGames := make([]domain.Game, 0, len(games))
	for _, g := range games {
		domainGames = append(domainGames, g.ToDomain())
	}

	return domainGames, nil
}

func buildSearchQuery(query string, fields, where string, limit, offset *int) string {
	apicalypse := fmt.Sprintf("search \"%s\"; fields %s;", query, fields)
	if where != "" {
		apicalypse += fmt.Sprintf(" where %s;", where)
	}
	if limit != nil {
		apicalypse += fmt.Sprintf(" limit %d;", *limit)
	}
	if offset != nil {
		apicalypse += fmt.Sprintf(" offset %d;", *offset)
	}
	return apicalypse
}
