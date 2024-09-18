package controller

import (
	"flash_sale_management/dto/request"
	"flash_sale_management/dto/response"
	"flash_sale_management/service"
	"flash_sale_management/utils"
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

// CreateFlashSale godoc
//
//	@Summary		Create Flash Sale
//	@Tags			Sales
//	@Accept			json
//	@Produce		json
//	@Param			request body request.CreateSaleRequest true "Request Body"
//	@Success		201 {object} response.SaleResponse "Created"
//	@Failure		400 {string} string "Bad Request"
//	@Router			/flash-sales [post]
func (s *SalesController) CreateFlashSale(c *fiber.Ctx) error {
	c.Accepts("application/json")
	saleRequest := new(request.CreateSaleRequest)

	if err := c.BodyParser(saleRequest); err != nil {
		return c.Status(http.StatusBadRequest).
			SendString(utils.CreateLogMessage("error parsing body", err))
	}

	sale, err := s.salesService.CreateSale(*saleRequest)
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	sale, err = s.salesService.SaveSale(sale)
	if err == nil {
		saleResponse := (&response.SaleResponse{}).FromEntity(sale)
		return c.Status(http.StatusCreated).JSON(saleResponse)
	} else {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}
}

// UpdateFlashSale ShowAccount godoc
//
//	@Summary		Update Flash Sale
//	@Tags			Sales
//	@Accept			json
//	@Produce		json
//	@Param			request body request.UpdateSaleRequest true "Request Body"
//	@Success		200 {object} response.SaleResponse "Ok"
//	@Failure		400 {string} string "Bad Request"
//	@Router			/flash-sales [put]
func (s *SalesController) UpdateFlashSale(c *fiber.Ctx) error {
	c.Accepts("application/json")
	saleRequest := new(request.UpdateSaleRequest)

	if err := c.BodyParser(saleRequest); err != nil {
		return c.Status(http.StatusBadRequest).
			SendString(utils.CreateLogMessage("error parsing body", err))
	}

	updatedSale, err := s.salesService.UpdateSale(*saleRequest)
	if err == nil {
		saleResponse := (&response.SaleResponse{}).FromEntity(updatedSale)
		return c.Status(http.StatusOK).JSON(saleResponse)
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
//	@Success		200 {object} response.SaleResponse "Ok"
//	@Failure		400 {string} string "Bad Request"
//	@Router			/flash-sales/{id} [get]
func (s *SalesController) GetFlashSale(c *fiber.Ctx) error {
	id := c.Params("id")
	saleID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(http.StatusBadRequest).
			SendString(utils.CreateLogMessage("wrong parameter. convert failed", err))
	}

	sales, err := s.salesService.FindSale(saleID)
	if err == nil {
		saleResponse := (&response.SaleResponse{}).FromEntity(sales)
		return c.Status(http.StatusOK).JSON(saleResponse)
	} else {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}
}

// GetFlashSales ShowAccount godoc
//
//	@Summary		Get Flash All Sale
//	@Tags			Sales
//	@Produce		json
//	@Success		200 {object} []response.SaleResponse "Ok"
//	@Failure		400 {string} string "Bad Request"
//	@Router			/flash-sales [get]
func (s *SalesController) GetFlashSales(c *fiber.Ctx) error {
	sales, err := s.salesService.FindSales()
	if err != nil {
		return c.Status(http.StatusBadRequest).
			SendString(utils.CreateLogMessage("error getting all sales", err))
	}

	var saleResponses []response.SaleResponse
	for _, sale := range *sales {
		if sale.ID == 0 {
			continue
		}

		saleResponse := (&response.SaleResponse{}).FromEntity(&sale)
		saleResponses = append(saleResponses, saleResponse)
	}

	if saleResponses == nil || len(saleResponses) == 0 {
		return c.Status(http.StatusOK).JSON("sales not found")
	}

	return c.Status(http.StatusOK).JSON(saleResponses)
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
		return c.Status(http.StatusBadRequest).
			SendString(utils.CreateLogMessage("wrong parameter. convert failed", err))
	}

	err = s.salesService.DeleteSale(saleID)
	if err != nil {
		log.Errorf(err.Error())
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
	// wait for transaction and race condition testing
	w8 := c.Query("wait", "1")

	id := c.Params("id")
	saleID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(http.StatusBadRequest).
			SendString(utils.CreateLogMessage("wrong parameter. convert failed", err))
	}

	wait, err := strconv.Atoi(w8)
	if err != nil {
		return c.Status(http.StatusBadRequest).
			SendString(utils.CreateLogMessage("wrong parameter. convert failed", err))
	}

	buy, err := s.salesService.Buy(saleID, wait)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(http.StatusOK).JSON(buy)
}
