package main

import "context"

func (sc *SearchClient) FuzzySearch(
	ctx context.Context,
	field, query string,
	fuzziness interface{},
) (*SearchResult, error) {
	// Fuzzy search for typo-tolerant searching. Match "laptop" even if typed as "latop". Uses Levenshtein distance
	// The Levenshtein distance (also known as edit distance) is a metric used to measure the difference between two strings. It calculates the minimum number of single-character operations required to transform one string into the other.
	searchQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"fuzzy": map[string]interface{}{
				field: map[string]interface{}{
					"value":     query,
					"fuzziness": fuzziness,
				},
			},
		},
	}

	return sc.executeSearch(ctx, searchQuery)
}
