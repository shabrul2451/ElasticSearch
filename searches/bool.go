package main

import "context"

func (sc *SearchClient) BoolSearch(
	ctx context.Context,
	params map[string]interface{},
) (*SearchResult, error) {
	// Boolean query with must, should, must_not
	searchQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": params,
		},
	}

	return sc.executeSearch(ctx, searchQuery)
}
