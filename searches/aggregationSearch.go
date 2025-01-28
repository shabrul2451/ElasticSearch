package main

import "context"

func (sc *SearchClient) AggregationSearch(
	ctx context.Context,
	aggs map[string]interface{},
) (*SearchResult, error) {
	// Search with aggregations. Performs statistical analysis. Can compute averages, sums, etc.
	searchQuery := map[string]interface{}{
		"size": 0, // We don't need hits for pure aggregations
		"aggs": aggs,
	}

	return sc.executeSearch(ctx, searchQuery)
}
