package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
)

// User represents a user document in Elasticsearch
type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// Config holds Elasticsearch configuration
type Config struct {
	Addresses []string
	Username  string
	Password  string
	APIKey    string
	Index     string
	CACert    string // Optional: Path to CA certificate
}

// ElasticsearchClient wraps the Elasticsearch client with basic CRUD operations
type ElasticsearchClient struct {
	client *elasticsearch.Client
	index  string
}

// NewElasticsearchClient creates a new Elasticsearch client
func NewElasticsearchClient(config Config) (*ElasticsearchClient, error) {
	cfg := elasticsearch.Config{
		Addresses: config.Addresses,
	}

	// Configure authentication - either use API key or username/password
	if config.APIKey != "" {
		cfg.APIKey = config.APIKey
	} else if config.Username != "" && config.Password != "" {
		cfg.Username = config.Username
		cfg.Password = config.Password
	} else {
		return nil, fmt.Errorf("either APIKey or Username/Password must be provided")
	}

	// Optional: Configure CA certificate if provided
	if config.CACert != "" {
		cfg.CACert = []byte(config.CACert)
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("error creating client: %w", err)
	}

	// Test the connection
	res, err := client.Info()
	if err != nil {
		return nil, fmt.Errorf("error connecting to Elasticsearch: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("error closing body")
		}
	}(res.Body)

	if res.IsError() {
		return nil, fmt.Errorf("error connecting to Elasticsearch: %s", res.String())
	}

	return &ElasticsearchClient{
		client: client,
		index:  config.Index,
	}, nil
}

func (c *ElasticsearchClient) CreateUser(user User) error {
	body, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("error marshalling user: %v", err)
	}

	res, err := c.client.Create(c.index, user.ID, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("error creating user: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("error closing body")
		}
	}(res.Body)
	if res.IsError() {
		return fmt.Errorf("error while creating user: %v", res)
	}
	return nil
}

func (c *ElasticsearchClient) UpdateUser(user User) error {
	body, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("error marshalling user: %v", err)
	}
	res, err := c.client.Update(c.index, user.ID, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("error updating user: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("error closing body")
		}
	}(res.Body)
	if res.IsError() {
		return fmt.Errorf("error while updating user: %v", res)
	}
	return nil
}

func (c *ElasticsearchClient) DeleteUser(user User) error {
	res, err := c.client.Delete(c.index, user.ID)
	if err != nil {
		return fmt.Errorf("error deleting user: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("error closing body")
		}
	}(res.Body)
	if res.IsError() {
		return fmt.Errorf("error while deleting user: %v", res)
	}
	return nil
}

func (c *ElasticsearchClient) GetUser(userId string) (*User, error) {
	res, err := c.client.Get(
		c.index,
		userId,
	)
	if err != nil {
		return nil, fmt.Errorf("error getting document: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("error closing body")
		}
	}(res.Body)

	if res.IsError() {
		if res.StatusCode == 404 {
			return nil, fmt.Errorf("document not found")
		}
		return nil, fmt.Errorf("error getting document: %s", res.String())
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}

	source := result["_source"].(map[string]interface{})
	user := &User{
		ID:    source["id"].(string),
		Name:  source["name"].(string),
		Email: source["email"].(string),
	}

	return user, nil
}

func (c *ElasticsearchClient) SearchUsers(query string) ([]User, error) {
	searchQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  query,
				"fields": []string{"name", "email"},
			},
		},
	}

	body, err := json.Marshal(searchQuery)
	if err != nil {
		return nil, fmt.Errorf("error marshalling search query: %v", err)
	}

	res, err := c.client.Search(
		c.client.Search.WithContext(context.Background()),
		c.client.Search.WithIndex(c.index),
		c.client.Search.WithBody(bytes.NewReader(body)),
	)
	if err != nil {
		return nil, fmt.Errorf("error searching users: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("error closing body")
		}
	}(res.Body)
	if res.IsError() {
		return nil, fmt.Errorf("error while searching users: %v", res)
	}

	var users []User
	if err := json.NewDecoder(res.Body).Decode(&users); err != nil {
		return nil, fmt.Errorf("error parsing search response: %v", err)
	}
	return users, nil
}

func main() {
	config := Config{
		Addresses: []string{"http://localhost:9200"},
		Username:  "elastic",  // Replace with your username
		Password:  "7FAW0rS2", // Replace with your password
		// OR use API key authentication
		// APIKey:    "your_api_key_here",
		Index: "users",
	}

	// Initialize Elasticsearch client
	ec, err := NewElasticsearchClient(config)
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	// Create a new user
	/*user := []User{
		{
			ID:        "3",
			Name:      "shabrul Doe",
			Email:     "john@example.com",
			CreatedAt: time.Now(),
		},
		{
			ID:        "2",
			Name:      "Jane Doe",
			Email:     "jane@example.com",
			CreatedAt: time.Now(),
		},
	}

	for _, user := range user {
		if err := ec.CreateUser(user); err != nil {
			log.Printf("Error creating user: %v", err)
		}
	}*/

	/*user, err := ec.GetUser("3")
	if err != nil {
		log.Fatalf("Error getting user: %v", err)
	}
	log.Printf("User: %v", user)*/

	if users, err := ec.SearchUsers("john"); err != nil {
		log.Printf("Error searching users: %v", err)
	} else {
		log.Printf("Found users: %+v", users)
	}
}
