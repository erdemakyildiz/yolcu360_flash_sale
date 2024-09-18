package response

import (
	"flash_sale_management/entity"
	"time"
)

type SaleResponse struct {
	ProductID int       `json:"product_id"`
	SaleStock int       `json:"saleStock"`
	Discount  float64   `json:"discount"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
	Active    bool      `json:"active"`
}

func (c *SaleResponse) FromEntity(sale *entity.Sale) SaleResponse {
	return SaleResponse{
		ProductID: sale.ProductID,
		SaleStock: sale.SaleStock,
		Discount:  sale.Discount,
		StartTime: sale.StartTime,
		EndTime:   sale.EndTime,
		Active:    sale.Active,
	}
}

type BuyProductResponse struct {
	ProductID             int       `json:"product_id"`
	RemainingSaleStock    int       `json:"RemainingSaleStock"`
	RemainingProductStock int       `json:"remainingProductStock"`
	Price                 float64   `json:"price"`
	BuyTime               time.Time `json:"time"`
}

func (c *BuyProductResponse) FromEntity(log entity.SaleLog) BuyProductResponse {
	return BuyProductResponse{
		ProductID:             log.ProductID,
		RemainingSaleStock:    log.RemainingSaleStock,
		RemainingProductStock: log.RemainingProductStock,
		Price:                 log.Price,
		BuyTime:               log.CreatedAt,
	}
}
