package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type productRepositoryRedis struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewProductRepositoryRedis(db *gorm.DB, redis *redis.Client) ProductRepository {
	db.AutoMigrate(&product{})
	mockData(db)
	return productRepositoryRedis{db: db, redis: redis}
}

func (r productRepositoryRedis) GetProducts() (products []product, err error) {
	key := "repository::GetProducts"

	// Redis Get
	productJson, err := r.redis.Get(context.Background(), key).Result()
	if err == nil {
		err = json.Unmarshal([]byte(productJson), &products)
		if err == nil {
			fmt.Println("redis")
			return products, nil
		}
	}

	// Database
	err = r.db.Order("quantity desc").Limit(30).Find(&products).Error
	if err != nil {
		return nil, err
	}

	// Redis set
	data, err := json.Marshal(products)
	if err != nil {
		return nil, err
	}

	err = r.redis.Set(context.Background(), key, string(data), time.Second*10).Err()
	if err != nil {
		return nil, err
	}

	fmt.Println("database")
	return products, nil
}
