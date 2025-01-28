package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/elastic/go-elasticsearch/v8"
)

type SearchParams struct {
	From int // Starting offset
	Size int // Number of results per page
}

type SearchClient struct {
	client *elasticsearch.Client
	index  string
}

// Product is a sample document structure
type Product struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       float64  `json:"price"`
	Categories  []string `json:"categories"`
	Brand       string   `json:"brand"`
	InStock     bool     `json:"in_stock"`
	Rating      float64  `json:"rating"`
}

// SearchResult represents the search response structure
type SearchResult struct {
	Total int64     `json:"total"`
	Items []Product `json:"items"`
	Aggs  any       `json:"aggregations,omitempty"`
}

func (sc *SearchClient) executeSearch(
	ctx context.Context,
	query map[string]interface{},
) (*SearchResult, error) {
	body, err := json.Marshal(query)
	if err != nil {
		return nil, fmt.Errorf("error marshaling query: %w", err)
	}

	res, err := sc.client.Search(
		sc.client.Search.WithContext(ctx),
		sc.client.Search.WithIndex(sc.index),
		sc.client.Search.WithBody(bytes.NewReader(body)),
	)
	if err != nil {
		return nil, fmt.Errorf("error executing search: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("error closing body")
		}
	}(res.Body)

	if res.IsError() {
		return nil, fmt.Errorf("search error: %s", res.String())
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}

	searchResult := &SearchResult{}

	// Extract total
	if hits, ok := result["hits"].(map[string]interface{}); ok {
		if total, ok := hits["total"].(map[string]interface{}); ok {
			searchResult.Total = int64(total["value"].(float64))
		}
	}

	// Extract items
	if hits, ok := result["hits"].(map[string]interface{}); ok {
		if hitsList, ok := hits["hits"].([]interface{}); ok {
			for _, hit := range hitsList {
				hitMap := hit.(map[string]interface{})
				source := hitMap["_source"].(map[string]interface{})

				var product Product
				sourceBytes, _ := json.Marshal(source)
				err := json.Unmarshal(sourceBytes, &product)
				if err != nil {
					return nil, err
				}
				searchResult.Items = append(searchResult.Items, product)
			}
		}
	}

	// Extract aggregations if present
	if aggs, ok := result["aggregations"].(map[string]interface{}); ok {
		searchResult.Aggs = aggs
	}

	return searchResult, nil
}

func toJson(res SearchResult) string {
	jsonData, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return ""
	}
	return string(jsonData)
}

func main() {
	ctx := context.Background()

	config := elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
		Username:  "elastic",
		Password:  "7FAW0rS2",
	}

	client, err := elasticsearch.NewClient(config)
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	sc := &SearchClient{
		client: client,
		index:  "products",
	}

	searchParams := SearchParams{
		From: 0,
		Size: 5,
	}

	// Example 1: Simple Match Search
	result, err := sc.MatchSearch(ctx, "name", "laptop", searchParams)
	if err != nil {
		log.Printf("Match search error: %v", err)
	}
	log.Println("match search result: ", toJson(*result))

	// Example 2: Multi-Match Search
	fields := []string{"name", "description"}
	result, err = sc.MultiMatchSearch(ctx, "gaming laptop", fields)
	if err != nil {
		log.Printf("Multi-match search error: %v", err)
	}
	log.Println("Multi-match search result: ", toJson(*result))

	// Example 3: Boolean Search[Find Apple products with price >= 1000]
	boolParams := map[string]interface{}{
		"must": []map[string]interface{}{
			{
				"match": map[string]interface{}{
					"brand": "Apple",
				},
			},
		},
		"filter": []map[string]interface{}{
			{
				"range": map[string]interface{}{
					"price": map[string]interface{}{
						"gte": 1000,
					},
				},
			},
		},
	}
	result, err = sc.BoolSearch(ctx, boolParams)
	if err != nil {
		log.Printf("Bool search error: %v", err)
	}
	log.Println("Boolean search result: ", toJson(*result))

	// Example 4: Range Search [Find products between $1000-$2000]
	ranges := map[string]interface{}{
		"gte": 1000,
		"lte": 2000,
	}
	result, err = sc.RangeSearch(ctx, "price", ranges)
	if err != nil {
		log.Printf("Range search error: %v", err)
	}
	log.Println("Range search result: ", toJson(*result))

	// Example 5: Fuzzy Search
	result, err = sc.FuzzySearch(ctx, "name", "lapto", 1)
	if err != nil {
		log.Printf("Fuzzy search error: %v", err)
	}
	log.Println("Fuzzy search result: ", toJson(*result))

	// Example 6: Aggregation Search
	aggs := map[string]interface{}{
		"avg_price": map[string]interface{}{
			"avg": map[string]interface{}{
				"field": "price",
			},
		},
		"categories": map[string]interface{}{
			"terms": map[string]interface{}{
				"field": "categories",
			},
		},
	}
	result, err = sc.AggregationSearch(ctx, aggs)
	if err != nil {
		log.Printf("Aggregation search error: %v", err)
	}
	log.Println("Aggregation search result: ", toJson(*result))

	// Example 7: Phrase Search
	result, err = sc.PhraseSearch(ctx, "description", "gaming laptop", 1)
	if err != nil {
		log.Printf("Phrase search error: %v", err)
	}
	log.Println("Phrase search result: ", toJson(*result))
}
