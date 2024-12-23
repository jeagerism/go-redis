package services

import (
	"context"
	"encoding/json"
	"fmt"
	"go-redis/repositories"
	"time"

	"github.com/go-redis/redis/v8"
)

type catalogServiceRedis struct {
	productRepo repositories.ProductRepository
	redisClient *redis.Client
}

func NewCatalogServiceRedis(productRepo repositories.ProductRepository, redisClient *redis.Client) CatalogService {
	return catalogServiceRedis{productRepo: productRepo, redisClient: redisClient}
}

func (s catalogServiceRedis) GetProducts() (products []Product, err error) {
	key := "service::GetProducts"

	// Redis GET
	if productJson, err := s.redisClient.Get(context.Background(), key).Result(); err == nil {
		fmt.Printf("Data from Redis: %s\n", productJson)
		if json.Unmarshal([]byte(productJson), &products) == nil {
			fmt.Printf("redis\n")
			return products, nil
		} else {
			fmt.Printf("Failed to unmarshal Redis data: %v\n", err)
		}
	} else {
		fmt.Printf("Failed to get data from Redis: %v\n", err)
	}

	// Repository
	productDB, err := s.productRepo.GetProducts()
	if err != nil {
		return nil, err
	}

	for _, p := range productDB {
		products = append(products, Product{
			ID:       p.ID,
			Name:     p.Name,
			Quantity: p.Quantity,
		})
	}

	// Redis Set
	if data, err := json.Marshal(products); err == nil {
		result := s.redisClient.Set(context.Background(), key, string(data), time.Minute)
		if result.Err() != nil {
			fmt.Printf("Failed to set data in Redis: %v\n", result.Err())
		}
	} else {
		fmt.Printf("Failed to marshal products: %v\n", err)
	}

	fmt.Println("database")
	return products, nil
}
