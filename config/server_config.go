package config

import (
	"context"
	"flash_sale_management/controller"
	"flash_sale_management/entity"
	"flash_sale_management/repository"
	"flash_sale_management/service"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"time"
)

func GetApplication() *fiber.App {

	// redis connection
	redisUri := viper.GetString("redis.connectionUri")
	client := redis.NewClient(&redis.Options{
		Addr:     redisUri,
		Password: "",
		DB:       0,
	})

	pong, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("Error connecting to Redis:", err)
	}
	fmt.Println("Connected to Redis:", pong)

	// postgres connection
	uri := viper.GetString("database.connectionUri")
	db, err := gorm.Open(postgres.Open(uri), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&entity.Product{}, &entity.Sale{}, &entity.SaleLog{})
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	debug := viper.GetBool("debug")
	if debug {
		db.Debug()
	}

	// redis service
	redisService := service.NewRedisService(client)

	// product service
	productRepository := repository.NewProductRepository(db)
	productService := service.NewProductService(productRepository, &redisService)

	// log service
	logRepository := repository.NewSaleLogRepository(db)
	logService := service.NewSaleLogService(logRepository)

	// sale service
	saleRepository := repository.NewSaleRepository(db)
	salesService := service.NewSalesService(saleRepository, productService, logService, &redisService)

	addTestProducts(productService)

	app := Handlers(controller.New(salesService))

	return app
}

func addTestProducts(productService service.ProductService) {
	//test product
	product := entity.Product{
		ID:        1,
		Name:      "Iphone 16",
		Price:     50.000,
		Stock:     10,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	product2 := entity.Product{
		ID:        2,
		Name:      "Iphone 17",
		Price:     100.000,
		Stock:     20,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	productService.CreateProduct(product)
	productService.CreateProduct(product2)
}
