package main

import "context"

func (sc *SearchClient) MatchSearch(
	ctx context.Context,
	field, query string,
	params SearchParams,
) (*SearchResult, error) {
	// Simple match query on a specific field. Good for basic text search.
	searchQuery := map[string]interface{}{
		"from": params.From,
		"size": params.Size,
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				field: query,
			},
		},
	}

	return sc.executeSearch(ctx, searchQuery)
}
