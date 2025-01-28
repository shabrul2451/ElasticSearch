package main

import "context"

func (sc *SearchClient) PhraseSearch(
	ctx context.Context,
	field, phrase string,
	slop int,
) (*SearchResult, error) {
	// Phrase search with slop. Slop in phrase search refers to the number of allowed word movements or rearrangements in a search query.Without Slop: The search engine requires an exact match of the words in the specified order.
	searchQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"match_phrase": map[string]interface{}{
				field: map[string]interface{}{
					"query": phrase,
					"slop":  slop,
				},
			},
		},
	}

	return sc.executeSearch(ctx, searchQuery)
}
