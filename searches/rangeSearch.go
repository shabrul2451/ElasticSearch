package main

import "context"

func (sc *SearchClient) RangeSearch(
	ctx context.Context,
	field string,
	ranges map[string]interface{},
) (*SearchResult, error) {
	// Range query for numeric/date fields
	searchQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"range": map[string]interface{}{
				field: ranges,
			},
		},
	}

	return sc.executeSearch(ctx, searchQuery)
}
