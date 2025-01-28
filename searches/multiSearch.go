package main

import "context"

func (sc *SearchClient) MultiMatchSearch(
	ctx context.Context,
	query string,
	fields []string,
) (*SearchResult, error) {
	// Search across multiple fields. Good for searching in title, description, etc.
	searchQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  query,
				"fields": fields,
			},
		},
	}

	return sc.executeSearch(ctx, searchQuery)
}
