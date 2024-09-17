package controller

import (
	"flash_sale_management/dto"
	"flash_sale_management/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"net/http"
	"strconv"
)

type SalesController struct {
	salesService service.SalesService
}

func New(salesService service.SalesService) SalesController {
	controller := SalesController{salesService: salesService}
	return controller
}

// CreateFlashSale ShowAccount godoc
//
//	@Summary		Create Flash Sale
//	@Tags			Sales
//	@Accept			json
//	@Produce		json
//	@Param			request body dto.CreateSaleRequest true "Request Body"
//	@Success  		201 "Created"
//	@Failure  		400 "Bad Request"
//	@Router			/flash-sales [post]
func (s *SalesController) CreateFlashSale(c *fiber.Ctx) error {
	c.Accepts("application/json")
	saleRequest := new(dto.CreateSaleRequest)

	if err := c.BodyParser(saleRequest); err != nil {
		log.Errorf("request validation error : %v", err.Error())
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}
	sale, err := s.salesService.CreateSale(*saleRequest)
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	_, err = s.salesService.SaveSale(sale)
	if err == nil {
		return c.SendStatus(http.StatusCreated)
	} else {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}
}

// GetFlashSale ShowAccount godoc
//
//	@Summary		Get Flash Sale
//	@Tags			Sales
//	@Produce		json
//	@Param			id path int true "Flash Sale ID"
//	@Success  		200 "Ok"
//	@Router			/flash-sales/{id} [get]
func (s *SalesController) GetFlashSale(c *fiber.Ctx) error {
	id := c.Params("id")
	saleID, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	sales, err := s.salesService.FindSale(saleID)
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	return c.JSON(sales)
}

// GetFlashSales ShowAccount godoc
//
//	@Summary		Get Flash All Sale
//	@Tags			Sales
//	@Produce		json
//	@Success  		200 "Ok"
//	@Router			/flash-sales [get]
func (s *SalesController) GetFlashSales(c *fiber.Ctx) error {
	sales, err := s.salesService.FindSales()
	if err != nil {
		return err
	}

	return c.JSON(sales)
}

// UpdateFlashSale ShowAccount godoc
//
//	@Summary		Update Flash Sale
//	@Tags			Sales
//	@Accept			json
//	@Produce		json
//	@Param			request body dto.UpdateSaleRequest true "Request Body"
//	@Success  		201 "Created"
//	@Failure  		400 "Bad Request"
//	@Router			/flash-sales [put]
func (s *SalesController) UpdateFlashSale(c *fiber.Ctx) error {
	c.Accepts("application/json")
	saleRequest := new(dto.UpdateSaleRequest)

	if err := c.BodyParser(saleRequest); err != nil {
		return err
	}
	_, err := s.salesService.UpdateSale(*saleRequest)
	if err == nil {
		return c.SendStatus(http.StatusOK)
	} else {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}
}

// DeleteFlashSale ShowAccount godoc
//
//	@Summary		Delete Flash Sale
//	@Tags			Sales
//	@Produce		json
//	@Param			id path int true "Flash Sale ID"
//	@Success  		200 "Ok"
//	@Router			/flash-sales/{id} [delete]
func (s *SalesController) DeleteFlashSale(c *fiber.Ctx) error {
	id := c.Params("id")
	saleID, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	_, err = s.salesService.DeleteSale(saleID)
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	return c.SendStatus(200)
}

// BuyProduct ShowAccount godoc
//
//	@Summary		Buy Product
//	@Tags			Sales
//	@Produce		json
//	@Success  		200 "Ok"
//	@Router			/flash-sales/{id}/buy [post]
func (s *SalesController) BuyProduct(c *fiber.Ctx) error {
	w8 := c.Query("wait", "5")
	id := c.Params("id")
	saleID, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	wait, err := strconv.Atoi(w8)
	if err != nil {
		return err
	}

	buy, err := s.salesService.Buy(saleID, wait)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(http.StatusOK).JSON(buy)
}
