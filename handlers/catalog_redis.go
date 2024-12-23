package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"go-redis/services"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

type catalogHandlerRedis struct {
	catalogSrv services.CatalogService
	redis      *redis.Client
}

func NewCatalogHandlerRedis(catalogSrv services.CatalogService, redis *redis.Client) CatalogHandler {
	return catalogHandlerRedis{catalogSrv: catalogSrv, redis: redis}
}

func (h catalogHandlerRedis) GetProducts(c *fiber.Ctx) error {
	key := "handler::GetProducts"

	// Redis Get
	if responseJson, err := h.redis.Get(context.Background(), key).Result(); err == nil {
		fmt.Println("redis")
		c.Set("Content-Type", "application/json")
		return c.SendString(responseJson)
	}

	// Services
	products, err := h.catalogSrv.GetProducts()
	if err != nil {
		return err
	}

	response := fiber.Map{
		"status":   "ok",
		"products": products,
	}

	// Redis SET
	if data, err := json.Marshal(response); err == nil {
		h.redis.Set(context.Background(), key, string(data), time.Minute)
	}

	fmt.Println("database")
	return c.JSON(response)
}
