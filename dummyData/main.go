package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
)

var ProductIndex = "products"

// Product represents a product in the catalog
type Product struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Categories  []string  `json:"categories"`
	Brand       string    `json:"brand"`
	InStock     bool      `json:"in_stock"`
	Rating      float64   `json:"rating"`
	CreatedAt   time.Time `json:"created_at"`
}

// Sample data arrays for generating random products
var (
	brands = []string{
		"Apple", "Samsung", "Dell", "HP", "Lenovo", "Asus", "Acer", "Microsoft",
		"LG", "Sony", "Intel", "AMD", "Razer", "MSI", "Toshiba",
	}

	categories = []string{
		"Laptops", "Smartphones", "Tablets", "Desktops", "Monitors",
		"Accessories", "Gaming", "Office", "Student", "Professional",
	}

	productTypes = []string{
		"Laptop", "Smartphone", "Tablet", "Desktop", "Monitor",
		"Keyboard", "Mouse", "Headphones", "Camera", "Printer",
	}

	adjectives = []string{
		"Professional", "Gaming", "Ultra", "Premium", "Basic",
		"Advanced", "Smart", "Portable", "Powerful", "Lightweight",
	}

	features = []string{
		"4K Display", "Touch Screen", "Fast Charging", "Wireless",
		"Bluetooth", "High Performance", "Long Battery Life",
		"Ergonomic Design", "RGB Lighting", "Compact",
	}
)

func generateProduct(id int) Product {
	rand.Seed(time.Now().UnixNano())

	// Generate random product name
	productType := productTypes[rand.Intn(len(productTypes))]
	adjective := adjectives[rand.Intn(len(adjectives))]
	brand := brands[rand.Intn(len(brands))]
	name := fmt.Sprintf("%s %s %s", brand, adjective, productType)

	// Generate random description
	feature1 := features[rand.Intn(len(features))]
	feature2 := features[rand.Intn(len(features))]
	description := fmt.Sprintf("A %s %s featuring %s and %s. Perfect for daily use.",
		adjective, productType, feature1, feature2)

	// Generate random categories (2-3 categories per product)
	numCategories := rand.Intn(2) + 2
	productCategories := make([]string, 0)
	for i := 0; i < numCategories; i++ {
		cat := categories[rand.Intn(len(categories))]
		if !contains(productCategories, cat) {
			productCategories = append(productCategories, cat)
		}
	}

	// Generate random price between 100 and 3000
	price := 100 + rand.Float64()*2900
	price = float64(int(price*100)) / 100 // Round to 2 decimal places

	return Product{
		ID:          fmt.Sprintf("%d", id),
		Name:        name,
		Description: description,
		Price:       price,
		Categories:  productCategories,
		Brand:       brand,
		InStock:     rand.Float32() > 0.2, // 80% chance of being in stock
		Rating:      1 + rand.Float64()*4, // Rating between 1 and 5
		CreatedAt:   time.Now().Add(-time.Duration(rand.Intn(365)) * 24 * time.Hour),
	}
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

type Mappings struct {
	Mappings struct {
		Properties struct {
			ID          Field `json:"id"`
			Name        Field `json:"name"`
			Description Field `json:"description"`
			Price       Field `json:"price"`
			Categories  Field `json:"categories"`
			Brand       Field `json:"brand"`
			InStock     Field `json:"in_stock"`
			Rating      Field `json:"rating"`
			CreatedAt   Field `json:"created_at"`
		} `json:"properties"`
	} `json:"mappings"`
}

type Field struct {
	Type string `json:"type"`
}

func seedData(client *elasticsearch.Client, indexName string, numProducts int) error {
	ctx := context.Background()

	// Delete index if exists
	_, err := client.Indices.Delete([]string{indexName})
	if err != nil {
		log.Printf("Error deleting index: %v", err)
	}

	// Create index with mappings
	mappings := Mappings{}
	mappings.Mappings.Properties.ID = Field{Type: "keyword"}
	mappings.Mappings.Properties.Name = Field{Type: "text"}
	mappings.Mappings.Properties.Description = Field{Type: "text"}
	mappings.Mappings.Properties.Price = Field{Type: "float"}
	mappings.Mappings.Properties.Categories = Field{Type: "keyword"}
	mappings.Mappings.Properties.Brand = Field{Type: "keyword"}
	mappings.Mappings.Properties.InStock = Field{Type: "boolean"}
	mappings.Mappings.Properties.Rating = Field{Type: "float"}
	mappings.Mappings.Properties.CreatedAt = Field{Type: "date"}

	jsonMappings, err := json.Marshal(mappings)
	if err != nil {
		return fmt.Errorf("error marshaling mappings: %w", err)
	}

	_, err = client.Indices.Create(
		indexName,
		client.Indices.Create.WithBody(bytes.NewReader(jsonMappings)),
	)
	if err != nil {
		return fmt.Errorf("error creating index: %w", err)
	}

	// Bulk indexing setup
	var bulk bytes.Buffer
	for i := 1; i <= numProducts; i++ {
		product := generateProduct(i)

		// Create bulk action line
		action := map[string]interface{}{
			"index": map[string]interface{}{
				"_id":    product.ID,
				"_index": indexName,
			},
		}

		actionJSON, err := json.Marshal(action)
		if err != nil {
			return fmt.Errorf("error marshaling action: %w", err)
		}
		bulk.Write(actionJSON)
		bulk.WriteString("\n")

		// Add product data
		productJSON, err := json.Marshal(product)
		if err != nil {
			return fmt.Errorf("error marshaling product: %w", err)
		}
		bulk.Write(productJSON)
		bulk.WriteString("\n")

		// Execute bulk request every 100 documents or on the last iteration
		if i%100 == 0 || i == numProducts {
			res, err := client.Bulk(
				bytes.NewReader(bulk.Bytes()),
				client.Bulk.WithContext(ctx),
				client.Bulk.WithIndex(indexName),
			)
			if err != nil {
				return fmt.Errorf("error bulk indexing: %w", err)
			}
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					fmt.Printf("error closing bulk indexing: %v", err)
				}
			}(res.Body)

			if res.IsError() {
				return fmt.Errorf("error bulk indexing: %s", res.String())
			}

			// Reset buffer for next batch
			bulk.Reset()
			log.Printf("Indexed %d products", i)
		}
	}

	// Refresh the index
	res, err := client.Indices.Refresh(
		client.Indices.Refresh.WithIndex(indexName),
	)
	if err != nil {
		return fmt.Errorf("error refreshing index: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("error closing bulk indexing: %v", err)
		}
	}(res.Body)

	if res.IsError() {
		return fmt.Errorf("error refreshing index: %s", res.String())
	}

	log.Printf("Successfully indexed %d products", numProducts)
	return nil
}

func main() {
	config := elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
		Username:  "elastic",
		Password:  "7FAW0rS2",
	}

	client, err := elasticsearch.NewClient(config)
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	// Seed 1000 products
	err = seedData(client, ProductIndex, 10000)
	if err != nil {
		log.Fatalf("Error seeding data: %v", err)
	}
}
