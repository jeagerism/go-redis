package main

import (
	"go-redis/handlers"
	"go-redis/repositories"
	"go-redis/services"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	db := initDatabase()
	redis := initRedis()
	_ = redis

	productRepo := repositories.NewProductRepositoryRedis(db, redis)
	productService := services.NewCatalogServiceRedis(productRepo, redis)
	productHandler := handlers.NewCatalogHandler(productService)

	app := fiber.New()
	app.Get("/products", productHandler.GetProducts)
	app.Listen(":8080")
}

func initDatabase() *gorm.DB {
	dial := mysql.Open("root:P@ssw0rd@tcp(localhost:3306)/infinitas")
	db, err := gorm.Open(dial, &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}

func initRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}
