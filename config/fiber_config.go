package config

import (
	"flash_sale_management/controller"
	_ "flash_sale_management/docs"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	"github.com/spf13/viper"
	"strconv"
)

func Handlers(controller controller.SalesController) *fiber.App {
	app := fiber.New()
	app.Use(cors.New())

	// sale
	app.Post("/flash-sales", controller.CreateFlashSale)
	app.Put("/flash-sales", controller.UpdateFlashSale)
	app.Get("/flash-sales", controller.GetFlashSales)
	app.Get("/flash-sales/:id", controller.GetFlashSale)
	app.Delete("/flash-sales/:id", controller.DeleteFlashSale)

	// buy product
	app.Post("/flash-sales/:id/buy", controller.BuyProduct)

	// swagger init
	app.Get("/swagger/*", swagger.HandlerDefault)

	return app
}

func StartServer() {
	LoadConfig()
	
	app := GetApplication()

	port := strconv.Itoa(viper.Get("server.port").(int))
	log.Fatal(app.Listen(":" + port))
}
