package config

import (
	"context"
	"flash_sale_management/controller"
	_ "flash_sale_management/docs"
	"flash_sale_management/entity"
	"flash_sale_management/repository"
	"flash_sale_management/service"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"strconv"
	"time"
)

func Handlers(controller controller.SalesController) *fiber.App {
	app := fiber.New()
	app.Use(cors.New())

	app.Post("/flash-sales", controller.CreateFlashSale)
	app.Get("/flash-sales", controller.GetFlashSales)
	app.Put("/flash-sales", controller.UpdateFlashSale)
	app.Get("/flash-sales/:id", controller.GetFlashSale)
	app.Delete("/flash-sales/:id", controller.DeleteFlashSale)
	app.Post("/flash-sales/:id/buy", controller.BuyProduct)

	app.Get("/swagger/*", swagger.HandlerDefault)

	return app
}

func StartServer() {
	LoadConfig()

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

	uri := viper.GetString("database.connectionUri")
	db, err := gorm.Open(postgres.Open(uri), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&entity.Product{}, &entity.Sale{}, &entity.SaleLog{})
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	db.Debug()

	productRepository := repository.NewProductRepository(db)
	productService := service.NewProductService(productRepository)

	logRepository := repository.NewSaleLogRepository(db)
	logService := service.NewSaleLogService(logRepository)

	redisService := service.NewRedisService(client)

	//test product
	product := entity.Product{
		ID:        20,
		Name:      "Test1",
		Price:     100,
		Stock:     10,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	product2 := entity.Product{
		ID:        21,
		Name:      "Test1",
		Price:     100,
		Stock:     0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	productService.CreateProduct(product)
	productService.CreateProduct(product2)

	saleRepository := repository.NewSaleRepository(db)
	salesService := service.NewSalesService(saleRepository, productService, logService, &redisService)

	app := Handlers(controller.New(salesService))
	port := strconv.Itoa(viper.Get("server.port").(int))
	log.Fatal(app.Listen(":" + port))
}
